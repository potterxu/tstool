package mpts

import (
	"fmt"
	"strings"

	"github.com/potterxu/bitreader"
)

type PrivateDataType struct {
	data []byte
}

type ExtensionType struct {
	data []byte
}

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
	PrivateData       *PrivateDataType

	ExtensionLength int
	Extension       *ExtensionType
}

func privateData(data []byte) *PrivateDataType {
	privateData := &PrivateDataType{
		data: data,
	}
	return privateData
}

func (privateData *PrivateDataType) String() string {
	var sb strings.Builder

	sb.WriteString("PrivateData: [ ")
	for _, v := range privateData.data {
		sb.WriteString(fmt.Sprintf("%x ", v))
	}
	sb.WriteString("]\n")

	return sb.String()
}

func extension(data []byte) *ExtensionType {
	extension := &ExtensionType{
		data: data,
	}
	return extension
}

func (extension *ExtensionType) String() string {
	var sb strings.Builder

	sb.WriteString("Extension: [ ")
	for _, v := range extension.data {
		sb.WriteString(fmt.Sprintf("%x ", v))
	}
	sb.WriteString("]\n")

	return sb.String()
}

func adaptation(data []byte) *AdaptationType {
	adapt := &AdaptationType{
		data: data,
	}
	adapt.parse()
	return adapt
}

func parsePCR(r *bitreader.BitReaderType) (int64, int64) {
	base, _ := r.ReadBits64(33)
	r.SkipBits(6)
	ext, _ := r.ReadBits64(9)
	return base, ext
}

func (adapt *AdaptationType) parse() {
	r := bitreader.BitReader(adapt.data)

	adapt.DI, _ = r.ReadBit()
	adapt.RAI, _ = r.ReadBit()
	adapt.ESPI, _ = r.ReadBit()
	adapt.PCRFlag, _ = r.ReadBit()
	adapt.OPCRFlag, _ = r.ReadBit()
	adapt.SplicingPointFlag, _ = r.ReadBit()
	adapt.PrivateDataFlag, _ = r.ReadBit()
	adapt.ExtensionFlag, _ = r.ReadBit()

	offset := 1

	if adapt.PCRFlag {
		base, ext := parsePCR(r)
		adapt.PCR = base*300 + ext
		offset += 6
	}

	if adapt.OPCRFlag {
		base, ext := parsePCR(r)
		adapt.OPCR = base*300 + ext
		offset += 6
	}

	if adapt.SplicingPointFlag {
		adapt.SpliceCountdown, _ = r.ReadBits(8)
		offset += 1
	}

	if adapt.PrivateDataFlag {
		adapt.PrivateDataLength, _ = r.ReadBits(8)
		offset += 1

		adapt.PrivateData = privateData(adapt.data[offset : offset+adapt.PrivateDataLength])
		r.SkipBytes(adapt.PrivateDataLength)
		offset += adapt.PrivateDataLength
	}

	if adapt.ExtensionFlag {
		adapt.ExtensionLength, _ = r.ReadBits(8)
		offset += 1

		adapt.Extension = extension(adapt.data[offset : offset+adapt.ExtensionLength])
		r.SkipBytes(adapt.ExtensionLength)
		offset += adapt.ExtensionLength
	}

	// stuffing
	r.SkipBytes(len(adapt.data) - offset)
}

func (adapt *AdaptationType) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintln("DI: ", adapt.DI))
	sb.WriteString(fmt.Sprintln("RAI: ", adapt.RAI))
	sb.WriteString(fmt.Sprintln("ESPI: ", adapt.ESPI))

	if adapt.PCRFlag {
		sb.WriteString(fmt.Sprintln("PCR: ", adapt.PCR))
	}

	if adapt.OPCRFlag {
		sb.WriteString(fmt.Sprintln("OPCR: ", adapt.OPCR))
	}

	if adapt.SplicingPointFlag {
		sb.WriteString(fmt.Sprintln("SpliceCountdown: ", adapt.SpliceCountdown))
	}

	if adapt.PrivateDataFlag {
		sb.WriteString(fmt.Sprintln(adapt.PrivateData))
	}

	if adapt.ExtensionFlag {
		sb.WriteString(fmt.Sprintln(adapt.Extension))
	}

	return sb.String()
}
