package mpts

import "github.com/potterxu/tstool/bitreader"

type AdaptationType struct {
	data []byte

	DI                bool // Discontinuity Indicator
	RAI               bool
	ESPI              bool
	PCRFlag           bool
	OPCRFlag          bool
	SplicingPointFlag bool
	PrivateDataFlag   bool
	ExtensionFlag     bool

	//todo
	PCR             int64
	OPCR            int64
	SpliceCountdown int

	PrivateDataLength int
	// PrivateData pointer

	ExtensionLength int
	// Extension pointer
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
	adapt.RAI = r.ReadBit()
	adapt.ESPI = r.ReadBit()
	adapt.PCRFlag = r.ReadBit()
	adapt.OPCRFlag = r.ReadBit()
	adapt.SplicingPointFlag = r.ReadBit()
	adapt.PrivateDataFlag = r.ReadBit()
	adapt.ExtensionFlag = r.ReadBit()

	remainSize := len(adapt.data) - 1

	if adapt.PCRFlag {
		base, ext := r.ReadPCR()
		adapt.PCR = base*300 + ext
		remainSize -= 6
	}

	if adapt.OPCRFlag {
		base, ext := r.ReadPCR()
		adapt.OPCR = base*300 + ext
		remainSize -= 6
	}

	if adapt.SplicingPointFlag {
		adapt.SpliceCountdown = r.ReadBits(8)
		remainSize -= 1
	}

	if adapt.PrivateDataFlag {
		adapt.PrivateDataLength = r.ReadBits(8)
		remainSize -= 1

		r.SkipBytes(adapt.PrivateDataLength)
		remainSize -= adapt.PrivateDataLength
	}

	if adapt.ExtensionFlag {
		adapt.ExtensionLength = r.ReadBits(8)
		remainSize -= 1

		r.SkipBytes(adapt.ExtensionLength)
		remainSize -= adapt.ExtensionLength
	}

	// stuffing
	r.SkipBytes(remainSize)
}
