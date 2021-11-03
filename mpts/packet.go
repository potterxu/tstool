package mpts

import "log"

const (
	TS_PACKET_SIZE int  = 188
	SYNC_BYTE      byte = 0x47
)

type PacketType struct {
	data []byte

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
	pkt.data = make([]byte, TS_PACKET_SIZE)
	copy(pkt.data, data)

	pkt.parse()
	return pkt
}

func (pkt *PacketType) parse() {
	upper := int(pkt.data[1] & 0x1f)
	lower := int(pkt.data[2] & 0xff)
	pkt.PID = upper<<8 + lower
}
