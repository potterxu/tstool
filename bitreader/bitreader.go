package bitreader

import (
	"fmt"

	"github.com/potterxu/tstool/util"
)

const (
	BYTE int = 8
)

type BitReaderType struct {
	data []byte

	base   int
	offset int
}

func BitReader(data []byte) *BitReaderType {
	return &BitReaderType{
		data:   data,
		base:   0,
		offset: 0,
	}
}

func (reader *BitReaderType) CheckBits(cnt int) bool {
	remain := len(reader.data) - reader.base
	remain *= 8
	remain -= reader.offset

	if remain < cnt {
		panic(fmt.Sprintf("[Out of bound] %d bits required, but only %d remained.", cnt, remain))
	}
	return true
}

func (reader *BitReaderType) SkipBytes(cnt int) {
	reader.SkipBits(cnt * BYTE)
}

func (reader *BitReaderType) SkipBits(cnt int) {
	cnt += reader.offset
	reader.base += cnt / BYTE
	reader.offset += cnt % BYTE
}

func (reader *BitReaderType) readInByte(cnt int) int64 {
	readBits := cnt
	unReadbits := BYTE - reader.offset - readBits
	mask := (0xff >> reader.offset) ^ (0xff >> (reader.offset + readBits))

	rv := (int64(reader.data[reader.base]) & int64(mask)) >> int64(unReadbits)

	return rv
}

func (reader *BitReaderType) ReadBits64(cnt int) int64 {
	reader.CheckBits(cnt)
	rv := int64(0)
	for cnt > 0 {
		readBits := util.MinInt(cnt, BYTE-reader.offset)
		rv <<= int64(readBits)
		rv += reader.readInByte(readBits)

		cnt -= readBits
		if readBits+reader.offset == BYTE {
			reader.base++
			reader.offset = 0
		} else {
			reader.offset += readBits
		}
	}
	return rv
}

func (reader *BitReaderType) ReadBits(cnt int) int {
	return int(reader.ReadBits64(cnt))
}

func (reader *BitReaderType) ReadBit() bool {
	return reader.ReadBits(1) != 0
}

func (reader *BitReaderType) ReadPCR() (int64, int64) {
	base := reader.ReadBits64(33)
	reader.SkipBits(6)
	ext := reader.ReadBits64(9)
	return base, ext
}
