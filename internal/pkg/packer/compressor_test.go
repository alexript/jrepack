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
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/itchio/lzma"
)

func TestSimplecompress(T *testing.T) {
	filename := "../../../test/output/simplecompress.dat"
	fd, _ := filepath.Abs(filename)
	defer os.Remove(fd)
	output, err := openOutput(filename)

	inputFolder := `../../../test/testdata/simplefolder`
	_, _, err = readInputFolder(inputFolder)

	T.Logf("Output struct: %v", output)
	written := closeOutput()

	T.Logf("Output file size: %d", written)
	if err != nil {
		T.Fatal(err)
	}

	T.Logf("Offsets table: %v", common.GetOffsets())

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		T.Fatal(err)
	}

	var b bytes.Buffer
	r := lzma.NewReader(f)
	_, _ = io.CopyN(&b, r, int64(written))
	_ = r.Close()
	readed := len(b.Bytes())
	if readed != 12 {
		T.Errorf("Unexpected uncompressed data size %d", readed)
	}

}
