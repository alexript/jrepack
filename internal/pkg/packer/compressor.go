package packer

import (
	"errors"
	"os"
)

type Output struct {
	File *os.File
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
	o = &Output{
		File: output,
	}
	return o, nil

}

func compress(data []byte) {
	if o != nil {
		o.File.Write(data)
	}
}

func closeOutput() {
	if o != nil {
		o.File.Close()
		o = nil
	}
}
