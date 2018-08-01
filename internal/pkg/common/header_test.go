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
	h.Fold(0, "test")
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

	f1 := NewFolder("f1")
	f2 := NewFolder("f2")
	f3 := NewFolder("f3")
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
	h := NewHeader(1000)
	f1id := h.Fold(0, "f1")
	f2id := h.Fold(f1id, "f2")
	h.Fold(f1id, "f3")
	h.Fold(f2id, "f4")

	bytes := ToBinary(h)
	T.Logf("Bytes: %v", bytes)

	h2 := FromBinary(bytes)

	T.Logf("Readed header: %v", h2)
	if h2.Size != 1000 {
		T.Errorf("Unexpected data size %d", h2.Size)
	}

	if len(h2.Folders) != 4 {
		T.Errorf("Unexpected folders number %d", len(h2.Folders))
	}
}
