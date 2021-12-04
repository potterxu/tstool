package io

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/potterxu/tstool/util"
)

type UdpReaderType struct {
	ipAddr string
	port   int
	interf string

	handle       *pcap.Handle
	packetSource *gopacket.PacketSource
}

func validateIPv4Address(ipAddr string, port int) bool {
	if net.ParseIP(ipAddr) == nil {
		log.Fatal("Invalid ip adderss ", ipAddr)
		return false
	}
	if port < 0 || port > 65535 {
		log.Fatal("Invalid port ", port)
		return false
	}

	return true
}

func UdpReader(ip string, port int, interf string) (*UdpReaderType, error) {
	if !validateIPv4Address(ip, port) {
		return nil, fmt.Errorf("invalid IPv4 address")
	}
	return &UdpReaderType{
		ipAddr: ip,
		port:   port,
		interf: interf,
	}, nil
}

func (udpReader *UdpReaderType) Open() {
	var err error
	udpReader.handle, err = pcap.OpenLive(udpReader.interf, 1504, false, 100*time.Millisecond)
	util.ExitOnErr(err)

	udpReader.packetSource = gopacket.NewPacketSource(udpReader.handle, udpReader.handle.LinkType())
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
		return nil, false
	}

	ip, _ := iplayer.(*layers.IPv4)
	if udpReader.ipAddr != ip.DstIP.String() {
		return nil, false
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer == nil {
		return nil, false
	}

	udp, _ := udpLayer.(*layers.UDP)
	if udpReader.port != int(udp.DstPort) {
		return nil, false
	}

	appLayer := packet.ApplicationLayer()
	if appLayer == nil {
		return nil, false
	}

	return appLayer.Payload(), true
}

func (udpReader *UdpReaderType) Write([]byte) {
	log.Fatal("Multicast write not supported yet")
}
