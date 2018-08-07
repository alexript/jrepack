package unpacker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	common "github.com/alexript/jrepack/internal/pkg/common"
)

func UnPack(inputFile, outputFolder string) error {
	input, err := filepath.Abs(inputFile)
	if err != nil {
		return err
	}
	ifi, err := os.Stat(input)
	if err != nil {
		return err
	}
	if ifi.IsDir() {
		return errors.New("Input file is the folder")
	}

	output, err := filepath.Abs(outputFolder)
	if err != nil {
		return err
	}
	_, err = os.Stat(output)
	if err == nil || os.IsExist(err) {
		return errors.New("Output folder exists")
	}

	header, err := readArch(inputFile)

	if err != nil {
		common.RemoveDirReq(output)
		return errors.New(fmt.Sprintf("Unable to read compressed header: %v", err))
	}

	err = Decompress(header, inputFile, output)
	header = nil
	runtime.GC()
	if err != nil {
		common.RemoveDirReq(output)
		return errors.New(fmt.Sprintf("Unable to decompress header: %v", err))
	}

	return err
}
