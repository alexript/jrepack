package packer

import (
	"errors"
	"io"
	"os"

	"github.com/itchio/lzma"
)

type Output struct {
	File   *os.File
	Writer io.WriteCloser
}

var (
	o *Output
)

func openOutput(filename string) (*Output, error) {
	if o != nil {
		return nil, errors.New("Output already open")
	}
	output, err := os.Create(filename)
	if err != nil {
		o = nil
		return nil, err
	}

	w := lzma.NewWriterLevel(output, 8)

	o = &Output{
		File:   output,
		Writer: w,
	}
	return o, nil

}

func compress(data []byte) {
	if o != nil {
		o.Writer.Write(data)
	}
}

func closeOutput() {
	if o != nil {
		o.Writer.Close()
		o.File.Close()
		o = nil
	}
}
