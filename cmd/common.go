package cmd

import (
	"github.com/potterxu/tstool/io"
)

func commonWorker(reader io.IOReaderInterface, writer io.IOWriterInterface) {
	for {
		buf, ok := reader.Read()
		if !ok {
			break
		}
		if buf != nil {
			writer.Write(buf)
		}
	}
}
