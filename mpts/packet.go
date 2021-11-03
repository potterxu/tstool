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

	PID int
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
	r.SkipBits(3)
	pkt.PID = r.ReadBits(13)
}
