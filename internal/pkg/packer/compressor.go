package packer

import (
	"errors"
	"io"
	"os"

	"github.com/itchio/lzma"
)

type Output struct {
	File        *os.File
	Writer      io.WriteCloser
	TailPointer int
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
		File:        output,
		Writer:      w,
		TailPointer: 0,
	}
	return o, nil

}

func compress(data []byte) (int, int, error) {
	if o != nil {
		offset := o.TailPointer
		n, err := o.Writer.Write(data)
		if err != nil {
			o.TailPointer += n
		}
		return offset, n, err
	}
	return -1, 0, nil
}

func closeOutput() {
	if o != nil {
		o.Writer.Close()
		o.File.Close()
		o = nil
	}
}
