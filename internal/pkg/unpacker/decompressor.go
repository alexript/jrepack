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

package unpacker

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/alexript/jrepack/ui"
	"github.com/itchio/lzma"
)

const (
	uint32max = (1 << 32) - 1
)

// GetOutputPath will transform outputdir string into disk path + path inside of archive
func GetOutputPath(h *common.Header, outputdir string, parentid uint32) (p *string, archp *string, err error) {

	if parentid < 1 {
		return &outputdir, nil, nil
	}

	parent := h.Folders[parentid-1]

	if string(parent.Name) == "_root_" {
		return &outputdir, nil, nil
	}

	pdir, adir, err := GetOutputPath(h, outputdir, parent.Parent)
	if err != nil {
		return nil, nil, err
	}
	if pdir == nil {
		return nil, nil, fmt.Errorf("Unable to determine upfolder for %v", parent)
	}

	var dirname string
	var archdir string
	var archdirp *string
	if parent.Flags == common.FArchive {
		dirname = path.Join(*pdir, string(parent.Name))
		archdir = ""
		archdirp = &archdir
	} else {
		if adir == nil {
			dirname = path.Join(*pdir, string(parent.Name))
			archdirp = nil
		} else {
			dirname = *pdir
			archdir = path.Join(*adir, string(parent.Name))
			archdirp = &archdir
		}
	}

	return &dirname, archdirp, nil
}

var (
	// OpenedZipFiles is the map of already opened archive files.
	OpenedZipFiles map[string]*os.File

	// ZipWriters is the map of already opened zip writers
	ZipWriters map[string]*zip.Writer
)

func initOpenedZipFiles() {
	OpenedZipFiles = make(map[string]*os.File)
	ZipWriters = make(map[string]*zip.Writer)
}

func closeOpenedZipFiles() {
	for _, writer := range ZipWriters {
		_ = writer.Close()
	}
	for _, file := range OpenedZipFiles {
		_ = file.Close()
	}
}

func saveToArch(archpath string, filename string, b []byte, isfolder bool) error {
	var targetFile *os.File
	var zipWriter *zip.Writer
	_, err := os.Stat(archpath)
	if err != nil && os.IsNotExist(err) {
		targetFile, err = os.Create(archpath)
		zipWriter = zip.NewWriter(targetFile)
		OpenedZipFiles[archpath] = targetFile
		ZipWriters[archpath] = zipWriter
	} else {
		if err == nil {
			targetFile = OpenedZipFiles[archpath]
			zipWriter = ZipWriters[archpath]
		} else {
			return err
		}
	}
	if err != nil {
		return err
	}

	fl := uint64(0)
	if b != nil {
		fl = uint64(len(b))
	}

	fh := &zip.FileHeader{
		Name:               filename,
		UncompressedSize64: fl,
	}
	fh.SetModTime(time.Now())
	fh.SetMode(0666)
	if fh.UncompressedSize64 > uint32max {
		fh.UncompressedSize = uint32max
	} else {
		fh.UncompressedSize = uint32(fh.UncompressedSize64)
	}
	if isfolder {
		fh.Name += "/"
	} else {
		fh.Method = zip.Deflate
	}

	writer, err := zipWriter.CreateHeader(fh)
	if err != nil {
		return err
	}
	if b != nil && !isfolder {
		_, err = writer.Write(b)
	}

	return err
}

func writeFile(outputdir string, header *common.Header, file *common.FolderRecord, b []byte) error {
	diskpath, archpath, err := GetOutputPath(header, outputdir, file.Parent)
	if err != nil {
		return err
	}

	if diskpath == nil {
		return fmt.Errorf("Unable to determine folder for output. %v", file)
	}

	if archpath == nil {
		// simple file
		err := os.MkdirAll(*diskpath, 0777)
		if err != nil {
			return err
		}
		filename := path.Join(*diskpath, string(file.Name))
		if file.Flags == common.FFolder {
			err := os.MkdirAll(filename, 0777)
			if err != nil {
				return err
			}
		} else {
			f, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer f.Close()
			if b != nil {
				_, err = f.Write(b)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// file to archive
		archdir := path.Dir(*diskpath)
		err := os.MkdirAll(archdir, 0777)
		if err != nil {
			return err
		}
		filename := path.Join(*archpath, string(file.Name))
		err = saveToArch(*diskpath, filename, b, file.Flags == common.FFolder)
		if err != nil {
			return err
		}
	}
	return nil
}

// Decompress is the entry point for decompressing process.
func Decompress(header *common.Header, filename string, output string) error {

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	initOpenedZipFiles()
	defer closeOpenedZipFiles()

	foldersNum := len(header.Folders)
	readedFolders := 0

	for _, folder := range header.Folders {
		if (folder.Flags == common.FData || folder.Flags == common.FFolder) && folder.Data == uint32(0xFFFFFFFF) {
			readedFolders++
			ui.Current().Unpack(readedFolders, foldersNum)
			err = writeFile(output, header, &folder, nil)
			if err != nil {
				return err
			}
		}
	}

	needToRead := int64(header.Size)
	readed := int64(0)

	var b bytes.Buffer

	r := lzma.NewReader(f)

	for _, dataRecord := range header.Data {
		b.Reset()
		n, err := io.CopyN(&b, r, int64(dataRecord.Size))
		if err != nil {
			return err
		}
		readed += n

		for _, folder := range header.Folders {
			if folder.Flags == common.FData && folder.Data == uint32(dataRecord.Offset) {
				readedFolders++
				err = writeFile(output, header, &folder, b.Bytes())
				ui.Current().Unpack(readedFolders, foldersNum)
				if err != nil {
					return err
				}
			}
		}

	}
	_ = r.Close()
	b.Reset()

	runtime.GC()
	if readed != needToRead {
		return fmt.Errorf("Readed: %d, Expected: %d", readed, needToRead)
	}

	return nil
}
