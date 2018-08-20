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
Package unpacker is the package for uncompress functions
*/
package unpacker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/alexript/jrepack/ui"
)

// UnPack is the entry pint of the package
func UnPack(inputFile, outputFolder string) error {
	input, err := filepath.Abs(inputFile)
	if err != nil {
		return err
	}
	ifi, err := os.Stat(input)
	if err != nil {
		return err
	}
	if ifi.IsDir() {
		return errors.New("Input file is the folder")
	}

	output, err := filepath.Abs(outputFolder)
	if err != nil {
		return err
	}
	_, err = os.Stat(output)
	if err == nil || os.IsExist(err) {
		return errors.New("Output folder exists")
	}

	header, err := readArch(inputFile)

	if err != nil {
		_ = common.RemoveDirReq(output)
		return fmt.Errorf("Unable to read compressed header: %v", err)
	}

	err = Decompress(header, inputFile, output)
	header = nil
	runtime.GC()
	if err != nil {
		_ = common.RemoveDirReq(output)
		return fmt.Errorf("Unable to decompress header: %v", err)
	}

	if err == nil {
		ui.Current().OnEnd(ui.EvtUnpackDone)
	}

	return err
}
