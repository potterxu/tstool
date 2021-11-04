package mpts

import (
	"log"

	"github.com/potterxu/tstool/bitreader"
)

const (
	TS_PACKET_SIZE int  = 188
	SYNC_BYTE      byte = 0x47
)

type PacketType struct {
	data [188]byte // stores the raw data

	TEI  bool // Transport error indicator
	PUSI bool // Payload unit start indicator
	TP   bool // Transport priority
	PID  int
	TSC  int // Transport scrambling control
	AFC  int // Adaptation field control
	CC   int // Continuity counter

	// Adaptation
	AdaptationLength int
	Adaptation       *AdaptationType

	// Payload
	PayloadLength int
	Payload       *PayloadType
}

func Packet(data []byte) *PacketType {
	if len(data) != TS_PACKET_SIZE {
		log.Fatal("Invaid ts packet size ", len(data))
		return nil
	}
	if data[0] != SYNC_BYTE {
		log.Fatal("Invalid sync byte ", data[0])
		return nil
	}
	pkt := &PacketType{}
	copy(pkt.data[:], data)

	pkt.parse()
	return pkt
}

func (pkt *PacketType) parse() {
	r := bitreader.BitReader(pkt.data[:])
	r.SkipBytes(1)
	pkt.TEI = r.ReadBit()
	pkt.PUSI = r.ReadBit()
	pkt.TP = r.ReadBit()
	pkt.PID = r.ReadBits(13)
	pkt.TSC = r.ReadBits(2)
	pkt.AFC = r.ReadBits(2)
	pkt.CC = r.ReadBits(4)
	pkt.AdaptationLength = 0
	pkt.PayloadLength = 0

	// Adaptation
	if (pkt.AFC & 0x10) > 0 {
		pkt.AdaptationLength = r.ReadBits(8)
		pkt.Adaptation = adaptation(pkt.data[5:5+pkt.AdaptationLength], r)
	}

	// Payload
	if (pkt.AFC & 0x01) > 0 {
		pkt.PayloadLength = TS_PACKET_SIZE - 5 - pkt.AdaptationLength
		pkt.Payload = payload(pkt.data[5+pkt.AdaptationLength:])
	}
}
