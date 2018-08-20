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

/*
Package packer is for jre packager.
*/
package packer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/alexript/jrepack/ui"
	"github.com/itchio/lzma"
)

/*
Pack is the entry point for package process.
*/
func Pack(inputFolder, outputFile string, dumpheader bool) error {
	input, err := filepath.Abs(inputFolder)
	if err != nil {
		return err
	}

	output, err := filepath.Abs(outputFile)
	if err != nil {
		return err
	}

	ifi, err := os.Stat(input)
	if err != nil {
		return err
	}
	if !ifi.IsDir() {
		return fmt.Errorf("Input folder %s is not the folder", input)
	}

	_, err = os.Stat(output)
	if err == nil {
		return fmt.Errorf("Output file %s exists", output)
	}

	_, err = openOutput(output)

	if err != nil {
		closeOutput()
		return err
	}

	_, rootfolder, err := readInputFolder(input)

	dataSize := closeOutput()

	h := common.NewHeader(dataSize)
	offsets := common.GetOffsets()
	h.Marshal(rootfolder, offsets)
	rootfolder = nil
	offsets = nil
	runtime.GC()
	binHeader := common.ToBinary(h)

	if dumpheader {
		json, err := os.Create(output + ".header.json")
		if err != nil {
			return err
		}
		defer json.Close()
		_, err = json.Write([]byte(h.String()))
		if err != nil {
			return err
		}
	}

	h = nil
	runtime.GC()

	if dumpheader {
		dump, err := os.Create(output + ".header")
		if err != nil {
			return err
		}
		defer dump.Close()
		_, err = dump.Write(binHeader)
		if err != nil {
			return err
		}
	}

	var compressedHeader bytes.Buffer
	w := lzma.NewWriterLevel(&compressedHeader, 8)
	_, err = w.Write(binHeader)
	if err != nil {
		runtime.GC()
		return err
	}
	err = w.Close()
	if err != nil {
		runtime.GC()
		return err
	}
	binHeader = nil
	runtime.GC()

	f, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()
	defer runtime.GC()

	chb := compressedHeader.Bytes()
	defer compressedHeader.Reset()

	_, err = f.Write(chb)
	if err != nil {
		return err
	}
	l := uint32(len(chb))
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, l)
	_, err = f.Write(a)

	if err == nil {
		ui.Current().OnEnd(ui.EvtPackDone)
	}

	return err
}
