package io

import (
	"bufio"
	"io"
	"os"

	"github.com/potterxu/tstool/util"
)

const (
	DEFAULT_BUFFER_SIZE int = 188 * 1000
	READ_BYTES              = 188
)

type FileWriterType struct {
	filename   string
	bufferSize int

	file   *os.File
	writer *bufio.Writer
}

func FileWriter(filename string) *FileWriterType {
	return &FileWriterType{
		filename:   filename,
		bufferSize: DEFAULT_BUFFER_SIZE,
	}
}

func (w *FileWriterType) Open() error {
	var err error
	w.file, err = os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.writer = bufio.NewWriterSize(w.file, w.bufferSize)
	return nil
}

func (w *FileWriterType) Close() {
	w.writer.Flush()
	w.file.Close()
}

func (w *FileWriterType) Write(data []byte) {
	_, err := w.writer.Write(data)
	util.LogOnErr(err)
}

type FileReaderType struct {
	filename   string
	bufferSize int

	file   *os.File
	reader *bufio.Reader
}

func FileReader(filename string) *FileReaderType {
	return &FileReaderType{
		filename:   filename,
		bufferSize: DEFAULT_BUFFER_SIZE,
	}
}

func (w *FileReaderType) Open() error {
	var err error
	w.file, err = os.Open(w.filename)
	if err != nil {
		return err
	}
	w.reader = bufio.NewReaderSize(w.file, w.bufferSize)
	return nil
}

func (w *FileReaderType) Close() {
	w.file.Close()
}

func (w *FileReaderType) Read() ([]byte, bool) {
	data := make([]byte, READ_BYTES)
	n, err := io.ReadFull(w.reader, data)

	return data, err == nil && n == READ_BYTES
}
