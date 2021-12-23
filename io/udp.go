package io

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	CUSTOM_EOF_TIMEOUT int = 100 // 100ms
)

type UdpReaderType struct {
	ipAddr string
	port   int
	interf string

	trafficAvailable bool

	handle       *pcap.Handle
	packetSource *gopacket.PacketSource
}

func UdpReader(ip string, port int, interf string) *UdpReaderType {
	return &UdpReaderType{
		ipAddr:           ip,
		port:             port,
		interf:           interf,
		trafficAvailable: false,
	}
}

func (udpReader *UdpReaderType) validateIPv4Address() error {
	if net.ParseIP(udpReader.ipAddr) == nil {
		return fmt.Errorf("invalid ip adderss ", udpReader.ipAddr)
	}
	if udpReader.port < 0 || udpReader.port > 65535 {
		return fmt.Errorf("Invalid port ", udpReader.port)
	}

	return nil
}

func (udpReader *UdpReaderType) Open() error {
	err := udpReader.validateIPv4Address()
	if err != nil {
		return err
	}
	udpReader.handle, err = pcap.OpenLive(udpReader.interf, 1504, false, 100*time.Millisecond)
	if err != nil {
		return err
	}
	filter := fmt.Sprintf("udp and dst host %v and dst port %v", udpReader.ipAddr, udpReader.port)
	err = udpReader.handle.SetBPFFilter(filter)
	if err != nil {
		return err
	}

	udpReader.packetSource = gopacket.NewPacketSource(udpReader.handle, udpReader.handle.LinkType())
	return nil
}

func (udpReader *UdpReaderType) Close() {
	udpReader.handle.Close()
}

func (udpReader *UdpReaderType) Read() ([]byte, bool) {
	timer := time.NewTimer(time.Duration(CUSTOM_EOF_TIMEOUT) * time.Millisecond)
	defer timer.Stop()

	var packet gopacket.Packet = nil
	var ok bool = true

	select {
	case <-timer.C:
		if udpReader.trafficAvailable {
			log.Default().Println("No traffic")
			udpReader.trafficAvailable = false
		}
	case packet, ok = <-udpReader.packetSource.Packets():
		// continue to process packet
		if !udpReader.trafficAvailable {
			log.Default().Println("Traffic detected")
			udpReader.trafficAvailable = true
		}
	}

	if !ok {
		return nil, false
	}

	if packet == nil {
		return nil, true
	}

	iplayer := packet.Layer(layers.LayerTypeIPv4)
	if iplayer == nil {
		return nil, true
	}

	ip, _ := iplayer.(*layers.IPv4)
	if udpReader.ipAddr != ip.DstIP.String() {
		return nil, true
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return nil, true
	}

	udp, _ := udpLayer.(*layers.UDP)
	if udpReader.port != int(udp.DstPort) {
		return nil, true
	}

	appLayer := packet.ApplicationLayer()
	if appLayer == nil {
		return nil, true
	}

	return appLayer.Payload(), true
}

func (udpReader *UdpReaderType) Write([]byte) {
	log.Fatal("Multicast write not supported yet")
}
