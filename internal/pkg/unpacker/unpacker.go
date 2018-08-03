package unpacker

import (
	"errors"
	"os"
	"path/filepath"

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
	if err == nil {
		return errors.New("Output folder exists")
	}

	err = nil

	_, err = readArch(inputFile)

	if err != nil {
		common.RemoveDirReq(output)
		return err
	}

	return nil
}
