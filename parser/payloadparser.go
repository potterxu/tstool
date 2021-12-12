package parser

import "github.com/potterxu/mpts"

type PayloadParserType struct {
	pid     int
	payload []byte
}

func PayloadParser(pid int) *PayloadParserType {
	p := &PayloadParserType{
		pid:     pid,
		payload: nil,
	}
	return p
}

func (p *PayloadParserType) Parse(pkt *mpts.PacketType) ([]byte, bool) {
	var rv []byte = nil
	var ok bool = false

	if pkt.PID != p.pid {
		return nil, false
	}

	if pkt.PUSI {
		if p.payload != nil {
			rv = p.payload
			ok = true
		}
		p.payload = make([]byte, 0)
	}

	if p.payload != nil {
		p.payload = append(p.payload, pkt.Payload.Data...)
	}

	return rv, ok
}
