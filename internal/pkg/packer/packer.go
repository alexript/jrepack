package packer

import (
	"errors"
	"os"
	"path/filepath"

	common "github.com/alexript/jrepack/internal/pkg/common"
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
	h.Marshal(rootfolder, common.GetOffsets())

	f, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(common.ToBinary(h))
	return err
}
