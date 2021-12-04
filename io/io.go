package io

type IOReaderInterface interface {
	Open()
	Close()
	Read() ([]byte, bool)
}

type IOWriterInterface interface {
	Open()
	Close()
	Write([]byte)
}

type IOInterface interface {
	Open()
	Close()
	Read() ([]byte, bool)
	Write([]byte)
}
