package packer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/itchio/lzma"
)

func Pack(inputFolder, outputFile string) error {
	input, err := filepath.Abs(inputFolder)
	if err != nil {
		return err
	}

	output, err := filepath.Abs(outputFile)
	if err != nil {
		return err
	}

	ifi, err := os.Stat(input)
	if err != nil {
		return err
	}
	if !ifi.IsDir() {
		return errors.New("Input folder is not the folder")
	}

	_, err = os.Stat(output)
	if err == nil {
		return errors.New("Output file exists")
	}

	_, err = openOutput(output)

	if err != nil {
		closeOutput()
		return err
	}

	_, rootfolder, err := readInputFolder(input)

	dataSize := closeOutput()

	h := common.NewHeader(dataSize)
	offsets := common.GetOffsets()
	h.Marshal(rootfolder, offsets)
	rootfolder = nil
	offsets = nil
	runtime.GC()
	binHeader := common.ToBinary(h)
	h = nil
	runtime.GC()
	var compressedHeader bytes.Buffer
	w := lzma.NewWriterLevel(&compressedHeader, 8)
	w.Write(binHeader)
	w.Close()
	binHeader = nil
	runtime.GC()

	f, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()
	defer runtime.GC()

	chb := compressedHeader.Bytes()
	defer compressedHeader.Reset()

	_, err = f.Write(chb)
	if err != nil {
		return err
	}
	l := uint32(len(chb))
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, l)
	_, err = f.Write(a)
	return err
}
