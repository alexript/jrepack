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
	"os"
	"path/filepath"
	"testing"
)

/*
TestPack tests packer.Pack function. packer.Pack should not to be
executed for not existed input and existed output.
*/
func TestPack(T *testing.T) {
	filename := "../../../test/output/packtest2.dat"
	inputFolder := `../../../test/testdata/simplefolder`
	f, _ := filepath.Abs(filename)
	_ = os.Remove(f)

	err := Pack("somefolder", "test", false)
	if err == nil {
		T.Fatal("Accept not existed input folder")
	}
	T.Log(err)

	err = Pack(".", "../../../test/output/.do_not_remove", false)
	if err == nil {
		T.Fatal("Accept existed output file")
	}
	T.Log(err)

	defer os.Remove(f)

	err = Pack(inputFolder, filename, false)
	if err != nil {
		T.Fatal(err)
	}
}

/*
TestPackZerosizeFile tests packer.Pack for correct packaging of files with 0 length.
*/
func TestPackZerosizeFile(T *testing.T) {
	filename := "../../../test/output/packtest3.dat"
	inputFolder := `../../../test/testdata/zerosize`
	f, _ := filepath.Abs(filename)
	_ = os.Remove(f)

	defer os.Remove(f)

	err := Pack(inputFolder, filename, false)
	if err != nil {
		T.Fatal(err)
	}

}
