package mpts

import (
	"fmt"
	"log"
	"strings"

	"github.com/potterxu/bitreader"
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
	pkt.TEI, _ = r.ReadBit()
	pkt.PUSI, _ = r.ReadBit()
	pkt.TP, _ = r.ReadBit()
	pkt.PID, _ = r.ReadBits(13)
	pkt.TSC, _ = r.ReadBits(2)
	pkt.AFC, _ = r.ReadBits(2)
	pkt.CC, _ = r.ReadBits(4)
	pkt.AdaptationLength = 0
	pkt.PayloadLength = 0

	// Adaptation
	if (pkt.AFC & 0x10) > 0 {
		pkt.AdaptationLength, _ = r.ReadBits(8)
		pkt.Adaptation = adaptation(pkt.data[5 : 5+pkt.AdaptationLength])
	}

	// Payload
	if (pkt.AFC & 0x01) > 0 {
		pkt.PayloadLength = TS_PACKET_SIZE - 5 - pkt.AdaptationLength
		pkt.Payload = payload(pkt.data[5+pkt.AdaptationLength:])
	}
}

func (pkt *PacketType) String() string {
	var sb strings.Builder
	sb.WriteString("TS Packet:\n")
	sb.WriteString("data: [ ")
	for _, v := range pkt.data {
		sb.WriteString(fmt.Sprintf("%x ", v))
	}
	sb.WriteString("]\n")
	sb.WriteString(fmt.Sprintln("TEI: ", pkt.TEI))
	sb.WriteString(fmt.Sprintln("PUSI: ", pkt.PUSI))
	sb.WriteString(fmt.Sprintln("TP: ", pkt.TP))
	sb.WriteString(fmt.Sprintln("PID: ", pkt.PID))
	sb.WriteString(fmt.Sprintln("TSC: ", pkt.TSC))
	sb.WriteString(fmt.Sprintln("AFC: ", pkt.AFC))
	sb.WriteString(fmt.Sprintln("CC: ", pkt.CC))

	return sb.String()
}
