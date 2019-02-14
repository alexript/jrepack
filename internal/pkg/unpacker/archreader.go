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

package unpacker

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/itchio/lzma"
)

func readArch(filename string) (*common.Header, error) {
	runtime.GC()

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(absPath)

	if os.IsNotExist(err) {
		return nil, errors.New("Path " + absPath + " does not exists")

	}
	if fi.IsDir() {
		return nil, errors.New(absPath + " is a folder")
	}

	filesize := fi.Size()

	f, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to open archive file: %v", err)
	}
	defer f.Close()

	point, err := f.Seek(-4, 2) // 4 bytes from end
	if err != nil {
		return nil, fmt.Errorf("Unable to seek for header size. FileSize: %v, Current point: %v, Error: %v", filesize, point, err)
	}
	b2 := make([]byte, 4)
	_, err = f.Read(b2)
	if err != nil {
		return nil, fmt.Errorf("Unable to read header size. FileSize: %v, Current point: %v, Error: %v", filesize, point, err)
	}
	packedHeaderSize := common.Order.Uint32(b2)

	point, err = f.Seek(-(int64(packedHeaderSize) + 4), 2)
	if err != nil {
		return nil, fmt.Errorf("Unable to seek for header head. FileSize: %v, Current point: %v, Error: %v", filesize, point, err)
	}
	b2 = make([]byte, packedHeaderSize)
	_, err = f.Read(b2)
	if err != nil {
		return nil, fmt.Errorf("Unable to read header: %v", err)
	}

	runtime.GC()

	br := bytes.NewReader(b2)
	var b bytes.Buffer
	r := lzma.NewReader(br)
	io.Copy(&b, r)
	r.Close()
	uncompressedHeader := b.Bytes()

	header := common.FromBinary(uncompressedHeader)
	uncompressedHeader = nil
	b.Reset()
	runtime.GC()
	return header, nil
}
