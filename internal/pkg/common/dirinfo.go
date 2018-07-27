package common

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"path"
	"strconv"
	"strings"
)

type file struct {
	Name    string
	Size    int
	Hashsum []byte
}

type folder struct {
	Name       string
	Folders    []*folder
	Files      []*file
	Containers []*container
}

type containerType struct {
	Name      string
	Extension string
}

type container struct {
	Content folder
	Type    *containerType
}

const (
	zipExt = ".zip"
	jarExt = ".jar"
)

var (
	zip            = containerType{Name: "zip file", Extension: zipExt}
	jar            = containerType{Name: "jar file", Extension: jarExt}
	containerTypes = []containerType{zip, jar}
	dirinfo        = make(map[string][]*file)
)

func isContainer(filename string) (*containerType, bool) {
	ext := path.Ext(filename)
	for _, v := range containerTypes {
		if strings.EqualFold(ext, v.Extension) {
			return &v, true
		}
	}
	return nil, false
}

func clearDirinfo() {
	dirinfo = make(map[string][]*file)
}

func addFileToDirinfo(f *file) {
	key := hex.EncodeToString(f.Hashsum)

	if _, ok := dirinfo[key]; !ok {
		dirinfo[key] = make([]*file, 0)
	}
	dirinfo[key] = append(dirinfo[key], f)
}

func newFile(filename string, body []byte) file {
	l := len(body)
	bs := []byte(strconv.Itoa(l))

	h := sha256.New()
	h.Reset()
	h.Write(bs) // hash is not just sha256 of file, but sha256 of file size _and_ file data
	h.Write(body)

	f := file{
		Name:    filename,
		Size:    l,
		Hashsum: h.Sum(nil),
	}

	addFileToDirinfo(&f)

	return f
}

func newFolder(foldername string) folder {
	return folder{
		Name:       foldername,
		Folders:    make([]*folder, 0),
		Files:      make([]*file, 0),
		Containers: make([]*container, 0),
	}
}

func newContainer(name string) *container {
	t, ok := isContainer(name)
	if !ok {
		return nil
	}
	return &container{
		Content: newFolder(name),
		Type:    t,
	}
}

func addFileToFolder(fold *folder, f *file) error {
	if fold == nil {
		return errors.New("Folder is nil")
	}
	if f == nil {
		return errors.New("File is nil")
	}
	fold.Files = append(fold.Files, f)
	return nil
}

func addFolderToFolder(dest *folder, src *folder) error {
	if dest == nil {
		return errors.New("Destination folder is nil")
	}
	if src == nil {
		return errors.New("Source folder is nil")
	}
	dest.Folders = append(dest.Folders, src)
	return nil
}

func addContainerToFolder(dest *folder, src *container) error {
	if dest == nil {
		return errors.New("Destination folder is nil")
	}
	if src == nil {
		return errors.New("Container is nil")
	}
	dest.Containers = append(dest.Containers, src)
	return nil
}

func addFileToContainer(c *container, f *file) error {
	if c == nil {
		return errors.New("Container is nil")
	}
	if f == nil {
		return errors.New("File is nil")
	}
	c.Content.Files = append(c.Content.Files, f)

	return nil
}

func addFolderToContainer(c *container, src *folder) error {
	if c == nil {
		return errors.New("Container is nil")
	}
	if src == nil {
		return errors.New("Source folder is nil")
	}
	c.Content.Folders = append(c.Content.Folders, src)
	return nil
}

func addContainerToContainer(dest *container, src *container) error {
	if dest == nil {
		return errors.New("Destination container is nil")
	}
	if src == nil {
		return errors.New("Source container is nil")
	}
	dest.Content.Containers = append(dest.Content.Containers, src)

	return nil
}
