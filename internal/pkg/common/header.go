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
	"bytes"
	"crypto/hmac"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
)

const (
	// FFolder is the bit for folder record
	FFolder uint8 = 0

	// FArchive is the bit for archive represented folder
	FArchive uint8 = 1

	// FData is the bit for file record
	FData uint8 = 2
)

var (
	// Order is the used binary endian order
	Order = binary.BigEndian
)

func getFlags(f *Folder) uint8 {
	flags := FFolder
	if f.IsContainer {
		flags |= FArchive
	}

	return flags
}

// FolderRecord is the type for folder representation in archive header.
type FolderRecord struct {
	Parent     uint32 `json:"parentId"`
	Flags      uint8  `json:"flags"`
	Data       uint32 `json:"dataId"`
	Namelength uint8  `json:"namelen"`
	Name       []byte `json:"name"`
}

// FoldersHeader is the array of the FolderRecord objects
type FoldersHeader []FolderRecord

// DataRecord is the type for file representation in archive header.
type DataRecord struct {
	Offset uint32 `json:"offset"`
	Size   uint32 `json:"size"`
	Hash   []byte `json:"hash"`
}

// DataHeader is the array of the pointers to the DataRecord objects.
type DataHeader []*DataRecord

// Header is the structure of the archive header information.
type Header struct {
	Folders FoldersHeader `json:"folders"`
	Data    DataHeader    `json:"data"`
	Size    uint32        `json:"datasize"`
}

func (h Header) String() string {
	dump, _ := json.MarshalIndent(h, "", "   ")
	return fmt.Sprintf("\nHeader:\n%s\n", dump)
}

// NewHeader will create new header object
func NewHeader(packedSize uint32) *Header {
	h := Header{
		Folders: make(FoldersHeader, 0),
		Data:    make(DataHeader, 0),
		Size:    packedSize,
	}
	return &h
}

// Packable is the interface for objects, which can be packed.
type Packable interface {
	Pack(offset uint32, size uint32, hash []byte)
}

// FindDataOffset will find data offset by hash of the given File object.
func (h *Header) FindDataOffset(f *File) uint32 {
	if f.Size == 0 {
		return 0xFFFFFFFF
	}
	for _, dr := range h.Data {
		if hmac.Equal(dr.Hash, f.Hashsum) {
			dr.Size = uint32(f.Size)
			return uint32(dr.Offset)
		}
	}
	return 0
}

// Fold will create new Folder record in header and link it with parent folder by parentID
func (h *Header) Fold(parentID uint32, f *Folder) uint32 {
	nbytes := []byte(f.Name)
	nl := len(nbytes)

	l := len(h.Folders) + 1
	folderID := uint32(l)

	flags := getFlags(f)
	rec := FolderRecord{
		Parent:     parentID,
		Flags:      flags,
		Data:       0xFFFFFFFF,
		Namelength: uint8(nl),
		Name:       nbytes,
	}

	h.Folders = append(h.Folders, rec)

	for _, file := range f.Files {
		nbytes = []byte(file.Name)
		nl = len(nbytes)
		rec = FolderRecord{
			Parent:     folderID,
			Flags:      FData,
			Data:       h.FindDataOffset(file),
			Namelength: uint8(nl),
			Name:       nbytes,
		}
		h.Folders = append(h.Folders, rec)
	}

	return folderID
}

// Pack will add new DataRecord into header
func (h *Header) Pack(offset uint32, size uint32, hash []byte) {
	rec := &DataRecord{
		Offset: offset,
		Size:   size,
		Hash:   hash,
	}
	h.Data = append(h.Data, rec)
}

func marsh(h *Header, parentID uint32, folders []*Folder) {
	for _, f := range folders {
		id := h.Fold(parentID, f)
		marsh(h, id, f.Folders)
	}
}

// Marshal will serialize root folder and offests into headr object
func (h *Header) Marshal(folder *Folder, offsets *Offset) {
	for offset, hash := range *offsets {
		h.Pack(offset, 0, hash)
	}
	id := h.Fold(0, folder)
	marsh(h, id, folder.Folders)

}

// ToBinary will transform header into bytearray
func ToBinary(h *Header) []byte {

	buf := new(bytes.Buffer)
	for _, f := range h.Folders {
		binary.Write(buf, Order, f.Parent)
		binary.Write(buf, Order, f.Flags)
		binary.Write(buf, Order, f.Data)
		binary.Write(buf, Order, f.Namelength)
		binary.Write(buf, Order, f.Name)
	}

	for _, d := range h.Data {
		binary.Write(buf, Order, d.Offset)
		binary.Write(buf, Order, d.Size)
		binary.Write(buf, Order, d.Hash)
	}

	binary.Write(buf, Order, uint32(len(h.Folders)))
	binary.Write(buf, Order, h.Size)
	return buf.Bytes()
}

// FromBinary will parse bytes array into new Header object
func FromBinary(b []byte) *Header {

	l := len(b)

	dataSize := Order.Uint32(b[l-4:])

	h := NewHeader(dataSize)
	foldersNum := Order.Uint32(b[l-8 : l-4])
	h.Folders = make(FoldersHeader, foldersNum)
	offset := uint32(0)
	for i := uint32(0); i < foldersNum; i++ {
		parentID := Order.Uint32(b[offset : offset+4])
		flags := uint8(b[offset+4])
		dataID := Order.Uint32(b[offset+5 : offset+9])

		namelength := uint8(b[offset+9])
		name := b[offset+10 : offset+10+uint32(namelength)]
		offset += (10 + uint32(namelength))
		rec := FolderRecord{
			Parent:     parentID,
			Flags:      flags,
			Data:       dataID,
			Namelength: namelength,
			Name:       name,
		}

		h.Folders[i] = rec
	}

	dataNum := ((l - 8) - int(offset)) / 40
	h.Data = make(DataHeader, dataNum)
	i := offset
	x := 0
	for i < uint32(l-8) {

		h.Data[x] = &DataRecord{
			Offset: Order.Uint32(b[i : i+4]),
			Size:   Order.Uint32(b[i+4 : i+8]),
			Hash:   b[i+8 : i+40],
		}

		i += 40
		x++
	}

	sort.Slice(h.Data, func(i, j int) bool { return h.Data[i].Offset < h.Data[j].Offset })
	runtime.GC()
	return h
}
