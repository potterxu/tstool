package io

type IOReaderInterface interface {
	Open() error
	Close()
	Read() ([]byte, bool)
}

type IOWriterInterface interface {
	Open() error
	Close()
	Write([]byte)
}

type IOInterface interface {
	Open() error
	Close()
	Read() ([]byte, bool)
	Write([]byte)
}
