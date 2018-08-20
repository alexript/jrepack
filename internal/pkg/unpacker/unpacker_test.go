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
	"os"
	"path/filepath"
	"testing"

	"github.com/alexript/jrepack/internal/pkg/common"
	"github.com/alexript/jrepack/internal/pkg/packer"
)

const (
	filenameTest      string = "../../../test/output/packtest.dat"
	inputFolderTest   string = `../../../test/testdata/mixedfolder`
	outputDirRootTest string = "../../../test/output/unpacked"
	outputDirNameTest string = "mixedfolder"
)

func prepareTestData() error {
	dropTestData()
	err := packer.Pack(inputFolderTest, filenameTest, false)
	return err
}

func dropTestData() {
	f, _ := filepath.Abs(filenameTest)
	os.Remove(f)

	root, _ := filepath.Abs(outputDirRootTest)

	dirName := filepath.Join(root, outputDirNameTest)

	common.RemoveDirReq(dirName)
}

func TestUnpacker(T *testing.T) {
	err := prepareTestData()
	if err != nil {
		T.Fatal(err)
	}
	defer dropTestData()

	err = UnPack("not_existed", outputDirRootTest)
	if err == nil {
		T.Error("Accepted not existed archive")
	}

	err = UnPack(filenameTest, outputDirRootTest)
	if err == nil {
		T.Error("Accepted existed target folder")
	}

	root, _ := filepath.Abs(outputDirRootTest)
	dirName := filepath.Join(root, outputDirNameTest)
	err = UnPack(filenameTest, dirName)
	if err != nil {
		T.Fatal(err)
	}
}
