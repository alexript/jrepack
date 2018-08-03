package common

import (
	"bytes"
	"crypto/hmac"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

const (
	FArchive uint8 = 1
	FData    uint8 = 2
)

var (
	Order = binary.BigEndian
)

func getFlags(f *Folder) uint8 {
	flags := uint8(0)
	if f.IsContainer {
		flags |= FArchive
	}

	return flags
}

type FolderRecord struct {
	Parent     uint32 `json:parentId`
	Flags      uint8  `json:flags`
	Data       uint32 `json:dataId`
	Namelength uint8  `json:namelen`
	Name       []byte `json:name`
}

type FoldersHeader []FolderRecord

type DataRecord struct {
	Offset uint32 `json:offset`
	Size   uint32 `json:size`
	Hash   []byte `json:hash`
}

type DataHeader []*DataRecord

type Header struct {
	Folders FoldersHeader `json:folders`
	Data    DataHeader    `json:data`
	Size    uint32        `json:datasize`
}

func (h Header) String() string {
	dump, _ := json.MarshalIndent(h, "", "   ")
	return fmt.Sprintf("\nHeader:\n%s\n", dump)
}

func NewHeader(packedSize uint32) *Header {
	h := Header{
		Folders: make(FoldersHeader, 0),
		Data:    make(DataHeader, 0),
		Size:    packedSize,
	}
	return &h
}

type Packable interface {
	Pack(offset uint32, size uint32, hash []byte)
}

func (h *Header) FindDataOffset(f *File) uint32 {
	for i, dr := range h.Data {
		if hmac.Equal(dr.Hash, f.Hashsum) {
			dr.Size = uint32(f.Size)
			return uint32(i)
		}
	}
	return 0
}

func (h *Header) Fold(parentId uint32, f *Folder) uint32 {
	nbytes := []byte(f.Name)
	nl := len(nbytes)

	l := len(h.Folders) + 1
	folderId := uint32(l)

	flags := getFlags(f)
	rec := FolderRecord{
		Parent:     parentId,
		Flags:      flags,
		Data:       0,
		Namelength: uint8(nl),
		Name:       nbytes,
	}

	h.Folders = append(h.Folders, rec)

	for _, file := range f.Files {
		nbytes = []byte(file.Name)
		nl = len(nbytes)
		rec = FolderRecord{
			Parent:     folderId,
			Flags:      FData,
			Data:       h.FindDataOffset(file),
			Namelength: uint8(nl),
			Name:       nbytes,
		}
		h.Folders = append(h.Folders, rec)
	}

	return folderId
}

func (h *Header) Pack(offset uint32, size uint32, hash []byte) {
	rec := &DataRecord{
		Offset: offset,
		Size:   size,
		Hash:   hash,
	}
	h.Data = append(h.Data, rec)
}

func marsh(h *Header, parentId uint32, folders []*Folder) {
	for _, f := range folders {
		id := h.Fold(parentId, f)
		marsh(h, id, f.Folders)
	}
}

func (h *Header) Marshal(folder *Folder, offsets *Offset) {
	for offset, hash := range *offsets {
		h.Pack(offset, 0, hash)
	}
	id := h.Fold(0, folder)
	marsh(h, id, folder.Folders)

}

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

func FromBinary(b []byte) *Header {

	l := len(b)

	dataSize := Order.Uint32(b[l-4:])

	h := NewHeader(dataSize)
	foldersNum := Order.Uint32(b[l-8 : l-4])
	h.Folders = make(FoldersHeader, foldersNum)
	offset := uint32(0)
	for i := uint32(0); i < foldersNum; i++ {
		parentId := Order.Uint32(b[offset : offset+4])
		flags := uint8(b[offset+4])
		dataId := Order.Uint32(b[offset+5 : offset+9])

		namelength := uint8(b[offset+9])
		name := b[offset+10 : offset+10+uint32(namelength)]
		offset += (10 + uint32(namelength))
		rec := FolderRecord{
			Parent:     parentId,
			Flags:      flags,
			Data:       dataId,
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

	return h
}
