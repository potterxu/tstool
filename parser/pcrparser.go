package parser

import (
	"log"

	"github.com/potterxu/mpts"
	"github.com/potterxu/tstool/algo"
)

type PcrParserType struct {
	pid    int
	pcrs   []int64
	pcrPos []int64

	pos int64
}

func PcrParser(pid int) *PcrParserType {
	return &PcrParserType{
		pid:    pid,
		pcrs:   make([]int64, 0),
		pcrPos: make([]int64, 0),
		pos:    0,
	}
}

func (p *PcrParserType) Parse(pkt *mpts.PacketType) (int64, int64, bool) {
	var pcr int64 = -1
	var pos int64 = -1
	var ok bool = false
	if pkt.PID == p.pid && pkt.AdaptationLength > 0 && pkt.Adaptation.PCRFlag {
		p.pcrs = append(p.pcrs, pkt.Adaptation.PCR)
		p.pcrPos = append(p.pcrPos, p.pos)

		pcr = pkt.Adaptation.PCR
		pos = p.pos
		ok = true
	}
	p.pos++
	return pcr, pos, ok
}

func (p *PcrParserType) GetPcr(pos int64) int64 {
	index := algo.BinarySearchNoLargerThan(p.pcrPos, pos)
	if p.pcrPos[index] == pos {
		return p.pcrs[index]
	}

	if index >= len(p.pcrPos)-1 || index < 0 {
		log.Default().Println("Cannot interpolate pcr for pos", pos)
		return -1
	}
	pcr1 := p.pcrs[index]
	pcr2 := p.pcrs[index+1]
	pos1 := p.pcrPos[index]
	pos2 := p.pcrPos[index+1]
	return pcr1 + (pcr2-pcr1)*(pos-pos1)/(pos2-pos1)
}
