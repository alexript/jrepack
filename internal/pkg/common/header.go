package common

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type FolderRecord struct {
	Parent     uint64 `json:parentId`
	Namelength uint8  `json:namelen`
	Name       []byte `json:name`
}

type FoldersHeader []FolderRecord

type DataRecord struct {
	Offset uint64 `json:offset`
	Size   uint64 `json:size`
	Hash   []byte `json:hash`
}

type DataHeader []DataRecord

type Header struct {
	Folders FoldersHeader `json:folders`
	Data    DataHeader    `json:data`
	Size    uint64        `json:datasize`
}

func (h Header) String() string {
	dump, _ := json.MarshalIndent(h, "", "   ")
	return fmt.Sprintf("\nHeader:\n%s\n", dump)
}

func NewHeader(packedSize uint64) *Header {
	h := Header{
		Folders: make(FoldersHeader, 0),
		Data:    make(DataHeader, 0),
		Size:    packedSize,
	}
	return &h
}

type Foldable interface {
	Fold(parentId uint64, name string) uint64
}

type Serializable interface {
	Marshal(folder *Folder, offsets *Offset)
}

type Packable interface {
	Pack(offset uint64, size uint64, hash []byte)
}

func (h *Header) Fold(parentId uint64, name string) uint64 {
	nbytes := []byte(name)
	nl := len(nbytes)
	rec := FolderRecord{
		Parent:     parentId,
		Namelength: uint8(nl),
		Name:       nbytes,
	}
	l := len(h.Folders) + 1
	h.Folders = append(h.Folders, rec)
	return uint64(l)
}

func (h *Header) Pack(offset uint64, size uint64, hash []byte) {
	rec := DataRecord{
		Offset: offset,
		Size:   size,
		Hash:   hash,
	}
	h.Data = append(h.Data, rec)
}

func marsh(h *Header, parentId uint64, folders []*Folder) {
	for _, f := range folders {
		id := h.Fold(parentId, f.Name)
		marsh(h, id, f.Folders)
	}
}

func (h *Header) Marshal(folder *Folder, offsets *Offset) {
	id := h.Fold(0, folder.Name)
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
		binary.Write(buf, order, f.Namelength)
		binary.Write(buf, order, f.Name)
	}

	for _, d := range h.Data {
		binary.Write(buf, order, d.Offset)
		binary.Write(buf, order, d.Size)
		binary.Write(buf, order, d.Hash)
	}

	binary.Write(buf, order, uint64(len(h.Folders)))
	binary.Write(buf, order, h.Size)
	return buf.Bytes()
}

func FromBinary(b []byte) *Header {
	order := binary.BigEndian
	l := len(b)

	dataSize := order.Uint64(b[l-8:])

	h := NewHeader(dataSize)
	foldersNum := order.Uint64(b[l-16 : l-8])
	offset := uint64(0)
	for i := uint64(0); i < foldersNum; i++ {
		parentId := order.Uint64(b[offset : offset+8])
		namelength := uint8(b[offset+8])
		name := b[offset+9 : offset+9+uint64(namelength)]
		offset += (9 + uint64(namelength))
		rec := FolderRecord{
			Parent:     parentId,
			Namelength: namelength,
			Name:       name,
		}

		h.Folders = append(h.Folders, rec)
	}
	return h
}
