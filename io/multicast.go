package io

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/potterxu/tstool/mpts"
	"github.com/potterxu/tstool/util"
)

type MulticastType struct {
	ipAddr     string
	port       int
	interf     string
	readerChan chan []byte

	handle       *pcap.Handle
	packetSource *gopacket.PacketSource

	running bool
	mutex   sync.RWMutex
}

func validateMulticastAddress(ipAddr string, port int) bool {
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

func Multicast(ipAddr string, port int, interf string) (*MulticastType, error) {
	if !validateMulticastAddress(ipAddr, port) {
		return nil, errors.New("invalid multicast address")
	}
	multicast := &MulticastType{
		ipAddr:     ipAddr,
		port:       port,
		interf:     interf,
		readerChan: make(chan []byte, 188*10000),
	}
	return multicast, nil
}

func (multicast *MulticastType) readWorker() {
	for packet := range multicast.packetSource.Packets() {
		{
			multicast.mutex.Lock()
			if !multicast.running {
				break
			}
			multicast.mutex.Unlock()
		}
		iplayer := packet.Layer(layers.LayerTypeIPv4)
		if iplayer == nil {
			continue
		}
		ip, _ := iplayer.(*layers.IPv4)
		if multicast.ipAddr != ip.DstIP.String() {
			continue
		}

		udpLayer := packet.Layer(layers.LayerTypeUDP)
		if udpLayer == nil {
			continue
		}
		udp, _ := udpLayer.(*layers.UDP)
		if multicast.port != int(udp.DstPort) {
			continue
		}

		appLayer := packet.ApplicationLayer()

		if appLayer != nil {
			// multicast.readerChan <- appLayer.Payload()
			buf := appLayer.Payload()
			for i := 0; i+mpts.TS_PACKET_SIZE <= len(buf); i += mpts.TS_PACKET_SIZE {
				multicast.readerChan <- buf[i : i+mpts.TS_PACKET_SIZE]
			}
		}
	}
}

func (multicast *MulticastType) Open() {
	var err error
	multicast.handle, err = pcap.OpenLive(multicast.interf, 1504, false, 30*time.Second)
	util.ExitOnErr(err)

	multicast.packetSource = gopacket.NewPacketSource(multicast.handle, multicast.handle.LinkType())
	multicast.running = true

	go multicast.readWorker()
}

func (multicast *MulticastType) Close() {
	multicast.handle.Close()
	{
		multicast.mutex.Lock()
		multicast.running = false
		multicast.mutex.Unlock()
	}
}

func (multicast *MulticastType) Read() []byte {
	return <-multicast.readerChan
}

func (multicast *MulticastType) Write([]byte) {
	log.Fatal("Multicast write not supported yet")
}
