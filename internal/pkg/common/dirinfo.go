package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type File struct {
	Name    string `json:name`
	Size    int    `json:size`
	Hashsum []byte `json:hash`
}

type Folder struct {
	Name       string       `json:name`
	Folders    []*Folder    `json:folders`
	Files      []*File      `json:files`
	Containers []*Container `json:containers`
}

type containerType struct {
	Name      string `json:name`
	Extension string `json:ext`
}

type Container struct {
	Content Folder         `json:content`
	Type    *containerType `json:type`
}

const (
	zipExt = ".zip"
	jarExt = ".jar"
)

type Offset map[uint64][]byte
type Dirinfo map[string][]*File

func (di Dirinfo) String() string {
	dump, _ := json.MarshalIndent(di, "", "   ")
	return fmt.Sprintf("\nDirinfo:\n%s\n", dump)
}

func (f Folder) String() string {
	dump, _ := json.MarshalIndent(f, "", "   ")
	return fmt.Sprintf("\nFolder:\n%s\n", dump)
}

func (c Container) String() string {
	dump, _ := json.MarshalIndent(c, "", "   ")
	return fmt.Sprintf("\nContainer:\n%s\n", dump)
}

func (o Offset) String() string {
	dump, _ := json.MarshalIndent(o, "", "   ")
	return fmt.Sprintf("\nOffsets table:\n%s\n", dump)
}

var (
	zip            = containerType{Name: "zip file", Extension: zipExt}
	jar            = containerType{Name: "jar file", Extension: jarExt}
	containerTypes = []containerType{zip, jar}
	dirinfo        = make(Dirinfo)
	offsets        = make(Offset)
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
	offsets = make(Offset)
}

func GetDirinfo() *Dirinfo {
	return &dirinfo
}

func GetOffsets() *Offset {
	return &offsets
}

func SetOffset(offset uint64, hash []byte) {
	offsets[offset] = hash
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

func NewFile(filename string, body []byte) (*File, bool) {
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

	return &f, isNewHash
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

type Foldernode interface {
	HasFolder(name string) (*Folder, error)
}

func findFolder(folders []*Folder, name string) (*Folder, error) {
	for _, f := range folders {
		if f.Name == name {
			return f, nil
		}
	}

	return nil, nil
}

func (fold *Folder) HasFolder(name string) (*Folder, error) {
	if fold == nil {
		return nil, errors.New("nil folder while search for " + name)
	}

	if len(fold.Folders) == 0 {
		return nil, nil
	}

	return findFolder(fold.Folders, name)

}

func (c *Container) HasFolder(name string) (*Folder, error) {
	if c == nil {
		return nil, errors.New("nil folder while search for " + name)
	}

	if len(c.Content.Folders) == 0 {
		return nil, nil
	}
	return findFolder(c.Content.Folders, name)
}

func MkdirAll(fold *Folder, somepath string) (f *Folder, err error) {
	l := len(somepath)
	i := 0
	for i < l && os.IsPathSeparator(somepath[i]) {
		i++
	}
	if i >= l {
		return fold, nil
	}
	j := i

	for i < l && !os.IsPathSeparator(somepath[i]) {
		i++
	}
	name := somepath[j:i]

	found, err := fold.HasFolder(name)
	if err != nil {
		return nil, err
	}

	var childfolder *Folder

	if found != nil {
		childfolder = found
	} else {
		t := NewFolder(name)
		childfolder = &t
		err = AddFolderToFolder(fold, childfolder)
		if err != nil {
			return nil, err
		}
	}
	if i >= l {
		return childfolder, nil
	}
	return MkdirAll(childfolder, somepath[i:])
}

func AddFileToFolder(fold *Folder, f *File) error {
	if fold == nil {
		return errors.New("Folder is nil")
	}
	if f == nil {
		return errors.New("File is nil")
	}

	dirname := path.Dir(f.Name)

	target := fold

	if dirname != "." {
		t, err := MkdirAll(target, dirname)
		if err != nil {
			return err
		}
		target = t
		f.Name = path.Base(f.Name)
	}

	target.Files = append(target.Files, f)
	return nil
}

func AddFolderToFolder(dest *Folder, src *Folder) error {
	if dest == nil {
		return errors.New("Destination folder is nil")
	}
	if src == nil {
		return errors.New("Source folder is nil")
	}

	dirname := path.Dir(src.Name)

	target := dest

	if dirname != "." {
		t, err := MkdirAll(target, dirname)
		if err != nil {
			return err
		}
		target = t
		src.Name = path.Base(src.Name)
	}
	target.Folders = append(target.Folders, src)
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

	dirname := path.Dir(f.Name)

	target := &(c.Content)

	if dirname != "." {
		t, err := MkdirAll(target, dirname)
		if err != nil {
			return err
		}
		target = t
		f.Name = path.Base(f.Name)
	}

	target.Files = append(target.Files, f)

	return nil
}

func AddFolderToContainer(c *Container, src *Folder) error {
	if c == nil {
		return errors.New("Container is nil")
	}
	if src == nil {
		return errors.New("Source folder is nil")
	}

	dirname := path.Dir(src.Name)

	target := &(c.Content)

	if dirname != "." {
		t, err := MkdirAll(target, dirname)
		if err != nil {
			return err
		}
		target = t
		src.Name = path.Base(src.Name)
	} else {
		target.Folders = append(target.Folders, src)
	}

	return nil
}

func AddContainerToContainer(dest *Container, src *Container) error {
	if dest == nil {
		return errors.New("Destination container is nil")
	}
	if src == nil {
		return errors.New("Source container is nil")
	}

	dirname := path.Dir(src.Content.Name)

	target := &(dest.Content)

	if dirname != "." {
		t, err := MkdirAll(target, dirname)
		if err != nil {
			return err
		}
		target = t
		src.Content.Name = path.Base(src.Content.Name)
	}
	target.Containers = append(target.Containers, src)

	return nil
}
