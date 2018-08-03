package unpacker

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/itchio/lzma"
)

func readArch(filename string) (*common.Header, error) {
	runtime.GC()

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(absPath)

	if os.IsNotExist(err) {
		return nil, errors.New("Path " + absPath + " does not exists")

	}
	if fi.IsDir() {
		return nil, errors.New(absPath + " is a folder")
	}

	filesize := fi.Size()

	f, err := os.Open(absPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to open archive file: %v", err))
	}
	defer f.Close()

	point, err := f.Seek(-4, 2) // 4 bytes from end
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to seek for header size. FileSize: %v, Current point: %v, Error: %v", filesize, point, err))
	}
	b2 := make([]byte, 4)
	_, err = f.Read(b2)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to read header size. FileSize: %v, Current point: %v, Error: %v", filesize, point, err))
	}
	packedHeaderSize := common.Order.Uint32(b2)

	point, err = f.Seek(-(int64(packedHeaderSize) + 4), 2)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to seek for header head. FileSize: %v, Current point: %v, Error: %v", filesize, point, err))
	}
	b2 = make([]byte, packedHeaderSize)
	_, err = f.Read(b2)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to read header: %v", err))
	}

	runtime.GC()

	br := bytes.NewReader(b2)
	var b bytes.Buffer
	r := lzma.NewReader(br)
	io.Copy(&b, r)
	r.Close()
	uncompressedHeader := b.Bytes()

	header := common.FromBinary(uncompressedHeader)
	uncompressedHeader = nil
	b.Reset()
	runtime.GC()
	return header, nil
}
