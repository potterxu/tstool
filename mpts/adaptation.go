package mpts

import "github.com/potterxu/tstool/bitreader"

type AdaptationType struct {
	data []byte

	DI bool // Discontinuity Indicator
}

func adaptation(data []byte, r *bitreader.BitReaderType) *AdaptationType {
	adapt := &AdaptationType{
		data: data,
	}
	adapt.parse(r)
	return adapt
}

func (adapt *AdaptationType) parse(r *bitreader.BitReaderType) {
	adapt.DI = r.ReadBit()
}
