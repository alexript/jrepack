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

package common

import (
	"testing"
)

func TestNewHeader(T *testing.T) {
	h := NewHeader(1000)
	if len(h.Folders) != 0 {
		T.Errorf("Unexpected folders size %d", len(h.Folders))
	}
	if h.Size != 1000 {
		T.Errorf("Unexpected compressed data size %v", h.Size)
	}
}

func TestAddFolder(T *testing.T) {
	h := NewHeader(1000)
	f1 := NewFolder("test", false)
	h.Fold(0, &f1)
	if len(h.Folders) != 1 {
		T.Errorf("Unexpected folders size %d", len(h.Folders))
	}
	fr := h.Folders[0]
	if fr.Parent != 0 {
		T.Errorf("Unexpected parentId %d", fr.Parent)
	}
	if fr.Namelength != 4 {
		T.Errorf("Unexpected name length %d", fr.Namelength)
	}

	if string(fr.Name) != "test" {
		T.Errorf("Unexpected name %v", fr.Name)
	}
}

func TestMarshalHeader(T *testing.T) {

	f1 := NewFolder("f1", false)
	f2 := NewFolder("f2", false)
	f3 := NewFolder("f3", false)
	AddFolderToFolder(&f2, &f3)
	AddFolderToFolder(&f1, &f2)

	o := Offset{}

	h := NewHeader(10)
	h.Marshal(&f1, &o)
	T.Logf("Header: %v", h)
	if len(h.Folders) != 3 {
		T.Errorf("Unexpected folders number: %v", len(h.Folders))
	}
}

func TestToFromBinary(T *testing.T) {
	f1 := NewFolder("f1", false)
	f2 := NewFolder("f2", false)
	f3 := NewFolder("f3", false)
	f4 := NewFolder("f4", true)

	h := NewHeader(1000)
	f1id := h.Fold(0, &f1)
	f2id := h.Fold(f1id, &f2)
	h.Fold(f1id, &f3)
	h.Fold(f2id, &f4)

	bytes := ToBinary(h)
	T.Logf("Bytes: %v", bytes)

	h2 := FromBinary(bytes)

	T.Logf("Readed header: %v", h2)
	if h2.Size != 1000 {
		T.Errorf("Unexpected data size %d", h2.Size)
	}

	if len(h2.Folders) != 4 {
		T.Fatalf("Unexpected folders number %d", len(h2.Folders))
	}

	f := h2.Folders[0] // f1
	if f.Parent != 0 {
		T.Errorf("Unexpected parent for f1: %d", f.Parent)
	}
	if f.Flags != 0 {
		T.Errorf("Unexpected flags for f1: %b", f.Flags)
	}
	if f.Namelength != 2 {
		T.Errorf("Unexpected name length for f1: %d", f.Namelength)
	}
	n := string(f.Name)
	if n != "f1" {
		T.Errorf("Unexpected name for f1: %s", n)
	}

	f = h2.Folders[1] // f2
	if f.Parent != 1 {
		T.Errorf("Unexpected parent for f2: %d", f.Parent)
	}
	if f.Flags != 0 {
		T.Errorf("Unexpected flags for f2: %b", f.Flags)
	}
	if f.Namelength != 2 {
		T.Errorf("Unexpected name length for f2: %d", f.Namelength)
	}
	n = string(f.Name)
	if n != "f2" {
		T.Errorf("Unexpected name for f2: %s", n)
	}

	f = h2.Folders[2] // f3
	if f.Parent != 1 {
		T.Errorf("Unexpected parent for f3: %d", f.Parent)
	}
	if f.Flags != 0 {
		T.Errorf("Unexpected flags for f3: %b", f.Flags)
	}
	if f.Namelength != 2 {
		T.Errorf("Unexpected name length for f3: %d", f.Namelength)
	}
	n = string(f.Name)
	if n != "f3" {
		T.Errorf("Unexpected name for f3: %s", n)
	}

	f = h2.Folders[3] // f4
	if f.Parent != 2 {
		T.Errorf("Unexpected parent for f4: %d", f.Parent)
	}
	if f.Flags != 1 {
		T.Errorf("Unexpected flags for f4: %b", f.Flags)
	}
	if f.Namelength != 2 {
		T.Errorf("Unexpected name length for f4: %d", f.Namelength)
	}
	n = string(f.Name)
	if n != "f4" {
		T.Errorf("Unexpected name for f4: %s", n)
	}
}
