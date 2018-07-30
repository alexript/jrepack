package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"path"
	"strconv"
	"strings"
)

type File struct {
	Name    string
	Size    int
	Hashsum []byte
}

type Folder struct {
	Name       string
	Folders    []*Folder
	Files      []*File
	Containers []*Container
}

type containerType struct {
	Name      string
	Extension string
}

type Container struct {
	Content Folder
	Type    *containerType
}

const (
	zipExt = ".zip"
	jarExt = ".jar"
)

type Dirinfo map[string][]*File

var (
	zip            = containerType{Name: "zip file", Extension: zipExt}
	jar            = containerType{Name: "jar file", Extension: jarExt}
	containerTypes = []containerType{zip, jar}
	dirinfo        = make(Dirinfo)
)

func IsContainer(filename string) (*containerType, bool) {
	ext := path.Ext(filename)
	for _, v := range containerTypes {
		if strings.EqualFold(ext, v.Extension) {
			return &v, true
		}
	}
	return nil, false
}

func ClearDirinfo() {
	dirinfo = make(Dirinfo)
}

func GetDirinfo() *Dirinfo {
	return &dirinfo
}

func addFileToDirinfo(f *File) bool {
	key := hex.EncodeToString(f.Hashsum)

	isNewHash := false

	if _, ok := dirinfo[key]; !ok {
		dirinfo[key] = make([]*File, 0)
		isNewHash = true
	}
	dirinfo[key] = append(dirinfo[key], f)
	return isNewHash
}

func NewFile(filename string, body []byte) (File, bool) {
	l := len(body)
	bs := []byte(strconv.Itoa(l))

	h := sha256.New()
	h.Reset()
	h.Write(bs) // hash is not just sha256 of file, but sha256 of file size _and_ file data
	h.Write(body)

	f := File{
		Name:    filename,
		Size:    l,
		Hashsum: h.Sum(nil),
	}

	isNewHash := addFileToDirinfo(&f)

	return f, isNewHash
}

func NewFolder(foldername string) Folder {
	return Folder{
		Name:       foldername,
		Folders:    make([]*Folder, 0),
		Files:      make([]*File, 0),
		Containers: make([]*Container, 0),
	}
}

func NewContainer(name string) *Container {
	t, ok := IsContainer(name)
	if !ok {
		return nil
	}
	return &Container{
		Content: NewFolder(name),
		Type:    t,
	}
}

func AddFileToFolder(fold *Folder, f *File) error {
	if fold == nil {
		return errors.New("Folder is nil")
	}
	if f == nil {
		return errors.New("File is nil")
	}
	fold.Files = append(fold.Files, f)
	return nil
}

func AddFolderToFolder(dest *Folder, src *Folder) error {
	if dest == nil {
		return errors.New("Destination folder is nil")
	}
	if src == nil {
		return errors.New("Source folder is nil")
	}
	dest.Folders = append(dest.Folders, src)
	return nil
}

func AddContainerToFolder(dest *Folder, src *Container) error {
	if dest == nil {
		return errors.New("Destination folder is nil")
	}
	if src == nil {
		return errors.New("Container is nil")
	}
	dest.Containers = append(dest.Containers, src)
	return nil
}

func AddFileToContainer(c *Container, f *File) error {
	if c == nil {
		return errors.New("Container is nil")
	}
	if f == nil {
		return errors.New("File is nil")
	}
	c.Content.Files = append(c.Content.Files, f)

	return nil
}

func AddFolderToContainer(c *Container, src *Folder) error {
	if c == nil {
		return errors.New("Container is nil")
	}
	if src == nil {
		return errors.New("Source folder is nil")
	}
	c.Content.Folders = append(c.Content.Folders, src)
	return nil
}

func AddContainerToContainer(dest *Container, src *Container) error {
	if dest == nil {
		return errors.New("Destination container is nil")
	}
	if src == nil {
		return errors.New("Source container is nil")
	}
	dest.Content.Containers = append(dest.Content.Containers, src)

	return nil
}
