package io

type IOInterface interface {
	Open()
	Close()
	Read() []byte
	Write([]byte)
}
