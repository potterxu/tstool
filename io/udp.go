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

type UdpReaderType struct {
	ipAddr string
	port   int
	interf string

	handle       *pcap.Handle
	packetSource *gopacket.PacketSource
}

func UdpReader(ip string, port int, interf string) *UdpReaderType {
	return &UdpReaderType{
		ipAddr: ip,
		port:   port,
		interf: interf,
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

	udpReader.packetSource = gopacket.NewPacketSource(udpReader.handle, udpReader.handle.LinkType())
	return nil
}

func (udpReader *UdpReaderType) Close() {
	udpReader.handle.Close()
}

func (udpReader *UdpReaderType) Read() ([]byte, bool) {
	packet, ok := <-udpReader.packetSource.Packets()

	if !ok {
		return nil, false
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
