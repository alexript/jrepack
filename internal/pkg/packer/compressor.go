package packer

import (
	"errors"
	"io"
	"os"
	"runtime"

	"github.com/alexript/jrepack/ui"
	"github.com/itchio/lzma"
)

type Output struct {
	File   *os.File
	Writer io.WriteCloser
}

var (
	o           *Output
	writtensize uint32
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
	writtensize = 0

	return o, nil

}

func compress(data []byte) (uint32, int, error) {
	if o != nil {
		l := len(data)
		offset := writtensize
		n, err := o.Writer.Write(data)
		if err == nil {
			writtensize = writtensize + uint32(l)

			ui.Current().Compress(ui.Compressed{
				Len:   l,
				Total: writtensize,
			})
		}

		return offset, n, err
	}
	return 0, 0, nil
}

func closeOutput() uint32 {
	if o != nil {
		o.Writer.Close()

		o.File.Close()
		o = nil
		runtime.GC()
	}
	return writtensize
}
