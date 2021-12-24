package parser

import "github.com/potterxu/mpts"

type PayloadParserType struct {
	pid     int
	payload []byte

	pktPos       int64
	lastStartPos int64
}

func PayloadParser(pid int) *PayloadParserType {
	p := &PayloadParserType{
		pid:          pid,
		payload:      nil,
		pktPos:       0,
		lastStartPos: -1,
	}
	return p
}

func (p *PayloadParserType) Parse(pkt *mpts.PacketType) ([]byte, int64, bool) {
	var rv []byte = nil
	var pos = p.lastStartPos
	var ok bool = false

	if pkt.PID == p.pid {
		if pkt.PUSI {
			p.lastStartPos = p.pktPos
			if p.payload != nil {
				rv = p.payload
				ok = true
			}
			p.payload = make([]byte, 0)
		}

		if p.payload != nil {
			p.payload = append(p.payload, pkt.Payload.Data...)
		}
	}

	p.pktPos++

	return rv, pos, ok
}

func (p *PayloadParserType) ParseWithPos(pkt *mpts.PacketType, pktPos int64) ([]byte, int64, bool) {
	var rv []byte = nil
	var pos = p.lastStartPos
	var ok bool = false

	if pkt.PID == p.pid {
		if pkt.PUSI {
			p.lastStartPos = pktPos
			if p.payload != nil {
				rv = p.payload
				ok = true
			}
			p.payload = make([]byte, 0)
		}

		if p.payload != nil && pkt.Payload != nil {
			p.payload = append(p.payload, pkt.Payload.Data...)
		}
	}

	return rv, pos, ok
}
