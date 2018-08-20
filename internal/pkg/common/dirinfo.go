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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/alexript/jrepack/ui"
)

// File is the representation of the file entity.
// This type contains filename, filesize and hashsumm.
type File struct {
	Name    string `json:name`
	Size    int    `json:size`
	Hashsum []byte `json:hash`
}

// Folder is the representation of the disk folder OR archive.
type Folder struct {
	IsContainer bool      `json:isContainer`
	Name        string    `json:name`
	Folders     []*Folder `json:folders`
	Files       []*File   `json:files`
}

// ContainerType is the type of container
type ContainerType struct {
	Name      string `json:name`
	Extension string `json:ext`
}

const (
	zipExt = ".zip"
	jarExt = ".jar"
)

// Offset is the file hashes by 4 bytes of file offset in _uncompressed_ data array.
type Offset map[uint32][]byte

// Dirinfo is the hash to files map, used as basic structure for output file header.
type Dirinfo map[string][]*File

func (di Dirinfo) String() string {
	dump, _ := json.MarshalIndent(di, "", "   ")
	return fmt.Sprintf("\nDirinfo:\n%s\n", dump)
}

func (f Folder) String() string {
	dump, _ := json.MarshalIndent(f, "", "   ")
	return fmt.Sprintf("\nFolder:\n%s\n", dump)
}

func (o Offset) String() string {
	dump, _ := json.MarshalIndent(o, "", "   ")
	return fmt.Sprintf("\nOffsets table:\n%s\n", dump)
}

var (
	zip            = ContainerType{Name: "zip file", Extension: zipExt}
	jar            = ContainerType{Name: "jar file", Extension: jarExt}
	containerTypes = []ContainerType{zip, jar}
	dirinfo        = make(Dirinfo)
	offsets        = make(Offset)
)

// IsContainer check file name for .zip or .jar extensions
func IsContainer(filename string) (*ContainerType, bool) {

	ext := path.Ext(filename)
	for _, v := range containerTypes {
		if strings.EqualFold(ext, v.Extension) {
			return &v, true
		}
	}
	return nil, false
}

// ClearDirinfo is for Dirinfo and Offset maps reset.
func ClearDirinfo() {
	dirinfo = make(Dirinfo)
	offsets = make(Offset)
}

// GetDirinfo will return current state of the DirInfo object
func GetDirinfo() *Dirinfo {
	return &dirinfo
}

// GetOffsets will return current state of the Offsets object
func GetOffsets() *Offset {
	return &offsets
}

// SetOffset apply hash to the data offset value.
func SetOffset(offset uint32, hash []byte) {
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

// NewFile will create new File object
func NewFile(filename string, body []byte) (*File, bool) {
	l := len(body)
	bs := []byte(strconv.Itoa(l))

	h := sha256.New()
	h.Reset()
	h.Write(bs) // hash is not just sha256 of file, but sha256 of file size _and_ file data
	h.Write(body)

	hs := h.Sum(nil)
	f := File{
		Name:    filename,
		Size:    l,
		Hashsum: hs,
	}

	isNewHash := addFileToDirinfo(&f)

	ui.Current().Hashed(ui.Hash{
		File:      filename,
		Size:      l,
		Hash:      hs,
		IsNewHash: isNewHash,
	})

	return &f, isNewHash
}

// NewFolder will create new Folder object.
func NewFolder(foldername string, isContainer bool) Folder {

	i := len(foldername) - 1

	for i > 0 && os.IsPathSeparator(foldername[i]) {
		i--
	}

	fname := foldername[0 : i+1]

	ui.Current().NewFolder(ui.Folder{
		IsContainer: isContainer,
		Name:        fname,
	})

	return Folder{
		IsContainer: isContainer,
		Name:        fname,
		Folders:     make([]*Folder, 0),
		Files:       make([]*File, 0),
	}
}

// Foldernode is the interface, defined is the object has subfolder with the given name.
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

// HasFolder will search for subfolder name in the given folder.
func (fold *Folder) HasFolder(name string) (*Folder, error) {
	if fold == nil {
		return nil, errors.New("nil folder while search for " + name)
	}

	if len(fold.Folders) == 0 {
		return nil, nil
	}

	return findFolder(fold.Folders, name)

}

// MkdirAll will create subfolders (if required) in the given folder.
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
		t := NewFolder(name, false)
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

// AddFileToFolder will add File object to Folder object
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

// AddFolderToFolder will add subfolder Folder object into Folder object
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
