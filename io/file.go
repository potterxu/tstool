package io

import (
	"bufio"
	"os"

	"github.com/potterxu/tstool/util"
)

const (
	DEFAULT_BUFFER_SIZE int = 188 * 1000
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
