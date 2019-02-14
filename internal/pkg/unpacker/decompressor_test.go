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
	"testing"

	common "github.com/alexript/jrepack/internal/pkg/common"
)

func TestGetOutputPath(T *testing.T) {
	f1 := common.NewFolder("_root_", false)
	f2 := common.NewFolder("f2", false)
	f3 := common.NewFolder("f3", false)
	f4 := common.NewFolder("f4", false)

	h := common.NewHeader(1000)
	f1id := h.Fold(0, &f1)
	f2id := h.Fold(f1id, &f2)
	f3id := h.Fold(f2id, &f3)
	f4id := h.Fold(f3id, &f4)

	T.Logf("Header: %v", h)

	dirpath, archpath, err := GetOutputPath(h, "", f4id)
	if err != nil {
		T.Error(err)
	}

	if archpath != nil {
		T.Errorf("Non-nil archive path %s", *archpath)
	}

	if *dirpath != "f2/f3/f4" {
		T.Errorf("Wrong path %s", *dirpath)
	}

	dirpath, archpath, err = GetOutputPath(h, "", 0)
	if err != nil {
		T.Error(err)
	}

	if archpath != nil {
		T.Errorf("Non-nil archive path %s", *archpath)
	}

	if *dirpath != "" {
		T.Errorf("Wrong path %s", *dirpath)
	}

	dirpath, archpath, err = GetOutputPath(h, "", f1id)
	if err != nil {
		T.Error(err)
	}

	if archpath != nil {
		T.Errorf("Non-nil archive path %s", *archpath)
	}

	if *dirpath != "" {
		T.Errorf("Wrong path %s", *dirpath)
	}

	dirpath, archpath, err = GetOutputPath(h, "", f2id)
	if err != nil {
		T.Error(err)
	}

	if archpath != nil {
		T.Errorf("Non-nil archive path %s", *archpath)
	}

	if *dirpath != "f2" {
		T.Errorf("Wrong path %s", *dirpath)
	}
}

func TestGetOutputPathWithArchive(T *testing.T) {
	f1 := common.NewFolder("_root_", false)
	f2 := common.NewFolder("f2", false)
	f3 := common.NewFolder("f3", true)
	f4 := common.NewFolder("f4", false)

	h := common.NewHeader(1000)
	f1id := h.Fold(0, &f1)
	f2id := h.Fold(f1id, &f2)
	f3id := h.Fold(f2id, &f3)
	f4id := h.Fold(f3id, &f4)

	T.Logf("Header: %v", h)

	dirpath, archpath, err := GetOutputPath(h, "", f4id)
	if err != nil {
		T.Error(err)
	}

	if archpath == nil {
		T.Fatal("Nil archive path")
	}

	if *dirpath != "f2/f3" || *archpath != "f4" {
		T.Errorf("Wrong path %s:%s", *dirpath, *archpath)
	}

}
