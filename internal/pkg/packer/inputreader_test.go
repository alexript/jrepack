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
	"testing"
)

func TestReadNotexistedInputFolder(T *testing.T) {
	_, _, err := readInputFolder("./non_existed_folder_name")
	if err == nil {
		T.Error("Can read non existed input folder")
	}
}

func TestReadExistedInputFolder(T *testing.T) {
	_, _, err := readInputFolder("../packer")
	if err != nil {
		T.Error(err)
	}
}

func TestReadInvalidInputFolder(T *testing.T) {
	_, _, err := readInputFolder(`\\\something`) // this is invalid filename for windows
	if err == nil {
		T.Error("Can read invalid input folder")
	}
}

func TestReadInputNotFolder(T *testing.T) {
	_, _, err := readInputFolder("packer.go")
	if err == nil {
		T.Error("Can read file as input folder")
	}
}

func TestReadSimplefolder(T *testing.T) {

	inputFolder := `../../../test/testdata/simplefolder`
	dirinfo, rootFolder, err := readInputFolder(inputFolder)
	if err != nil || dirinfo == nil {
		T.Errorf("Unable to read test data from %v", inputFolder)
	}

	T.Logf("Dirinfo: %v", dirinfo)
	T.Logf("Root folder: %v", rootFolder)

	resultDirinfoLength := len(*dirinfo)
	if resultDirinfoLength != 2 {
		T.Errorf("Expected only 1 record in dirinfo. Result: %v", resultDirinfoLength)
	}
}

func TestReadSimplecontainer(T *testing.T) {

	inputFolder := `../../../test/testdata/simplecontainer`
	dirinfo, rootFolder, err := readInputFolder(inputFolder)
	if err != nil || dirinfo == nil {
		T.Errorf("Unable to read test data from %v", inputFolder)
	}

	T.Logf("Dirinfo: %v", dirinfo)
	T.Logf("Root folder: %v", rootFolder)

	resultDirinfoLength := len(*dirinfo)
	if resultDirinfoLength != 1 {
		T.Errorf("Expected only 1 record in dirinfo. Result: %v", resultDirinfoLength)
	}
}

func TestReadNestedcontainer(T *testing.T) {

	inputFolder := `../../../test/testdata/nestedcontainer`
	dirinfo, rootFolder, err := readInputFolder(inputFolder)
	if err != nil || dirinfo == nil {
		T.Errorf("Unable to read test data from %v", inputFolder)
	}

	T.Logf("Dirinfo: %v", dirinfo)
	T.Logf("Root folder: %v", rootFolder)

	resultDirinfoLength := len(*dirinfo)
	if resultDirinfoLength != 1 {
		T.Errorf("Expected only 1 record in dirinfo. Result: %v", resultDirinfoLength)
	}
}
