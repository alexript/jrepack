package common

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

const (
	FArchive uint8 = 1
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
	Namelength uint8  `json:namelen`
	Name       []byte `json:name`
}

type FoldersHeader []FolderRecord

type DataRecord struct {
	Offset uint32 `json:offset`
	Size   uint32 `json:size`
	Hash   []byte `json:hash`
}

type DataHeader []DataRecord

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

type Foldable interface {
	Fold(parentId uint32, f *Folder) uint32
}

type Serializable interface {
	Marshal(folder *Folder, offsets *Offset)
}

type Packable interface {
	Pack(offset uint32, size uint32, hash []byte)
}

func (h *Header) Fold(parentId uint32, f *Folder) uint32 {
	nbytes := []byte(f.Name)
	nl := len(nbytes)
	flags := getFlags(f)
	rec := FolderRecord{
		Parent:     parentId,
		Flags:      flags,
		Namelength: uint8(nl),
		Name:       nbytes,
	}
	l := len(h.Folders) + 1
	h.Folders = append(h.Folders, rec)
	return uint32(l)
}

func (h *Header) Pack(offset uint32, size uint32, hash []byte) {
	rec := DataRecord{
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
	id := h.Fold(0, folder)
	marsh(h, id, folder.Folders)

	for offset, hash := range *offsets {
		h.Pack(offset, 0, hash)
	}
}

func ToBinary(h *Header) []byte {
	order := binary.BigEndian
	buf := new(bytes.Buffer)
	for _, f := range h.Folders {
		binary.Write(buf, order, f.Parent)
		binary.Write(buf, order, f.Flags)
		binary.Write(buf, order, f.Namelength)
		binary.Write(buf, order, f.Name)
	}

	for _, d := range h.Data {
		binary.Write(buf, order, d.Offset)
		binary.Write(buf, order, d.Size)
		binary.Write(buf, order, d.Hash)
	}

	binary.Write(buf, order, uint32(len(h.Folders)))
	binary.Write(buf, order, h.Size)
	return buf.Bytes()
}

func FromBinary(b []byte) *Header {
	order := binary.BigEndian
	l := len(b)

	dataSize := order.Uint32(b[l-4:])

	h := NewHeader(dataSize)
	foldersNum := order.Uint32(b[l-8 : l-4])
	offset := uint32(0)
	for i := uint32(0); i < foldersNum; i++ {
		parentId := order.Uint32(b[offset : offset+4])
		flags := uint8(b[offset+4])
		namelength := uint8(b[offset+5])
		name := b[offset+6 : offset+6+uint32(namelength)]
		offset += (6 + uint32(namelength))
		rec := FolderRecord{
			Parent:     parentId,
			Flags:      flags,
			Namelength: namelength,
			Name:       name,
		}

		h.Folders = append(h.Folders, rec)
	}
	return h
}
