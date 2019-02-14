// Copyright (C) 2018  Alexander Malyshev

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package packer

import (
	"errors"
	"io"
	"os"
	"runtime"

	"github.com/alexript/jrepack/ui"
	"github.com/itchio/lzma"
)

// Output is container for lzma writer object
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
		_ = o.Writer.Close()

		_ = o.File.Close()
		o = nil
		runtime.GC()
	}
	return writtensize
}
