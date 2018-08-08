package unpacker

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/itchio/lzma"
)

const (
	uint32max = (1 << 32) - 1
)

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
		return nil, nil, errors.New(fmt.Sprintf("Unable to determine upfolder for %v", parent))
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
	OpenedZipFiles map[string]*os.File
	ZipWriters     map[string]*zip.Writer
)

func initOpenedZipFiles() {
	OpenedZipFiles = make(map[string]*os.File)
	ZipWriters = make(map[string]*zip.Writer)
}

func closeOpenedZipFiles() {
	for _, writer := range ZipWriters {
		writer.Close()
	}
	for _, file := range OpenedZipFiles {
		file.Close()
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
		return errors.New(fmt.Sprintf("Unable to determine folder for output. %v", file))
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

func Decompress(header *common.Header, filename string, output string) error {

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}

	initOpenedZipFiles()
	defer closeOpenedZipFiles()

	for _, folder := range header.Folders {
		if (folder.Flags == common.FData || folder.Flags == common.FFolder) && folder.Data == uint32(0xFFFFFFFF) {
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
				err = writeFile(output, header, &folder, b.Bytes())
				if err != nil {
					return err
				}
			}
		}

	}
	r.Close()
	b.Reset()

	runtime.GC()
	if readed != needToRead {
		return errors.New(fmt.Sprintf("Readed: %d, Expected: %d", readed, needToRead))
	}

	return nil
}
