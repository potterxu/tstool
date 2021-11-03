package bitreader

import (
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

func (reader *BitReaderType) SkipBytes(cnt int) {
	reader.base += cnt
}

func (reader *BitReaderType) SkipBits(cnt int) {
	cnt += reader.offset
	reader.base += cnt / BYTE
	reader.offset += cnt % BYTE
}

func (reader *BitReaderType) readInByte(cnt int) int64 {
	readBits := cnt
	unReadbits := BYTE - reader.offset - readBits
	mask := (1<<readBits - 1) ^ (1<<unReadbits - 1)

	rv := (int64(reader.data[reader.base]) & int64(mask)) >> int64(unReadbits)

	return rv
}

func (reader *BitReaderType) ReadBits64(cnt int) int64 {
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
