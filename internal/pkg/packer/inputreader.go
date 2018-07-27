package packer

import (
	"errors"
	"os"
	"path/filepath"
)

func readInputFolder(inputFolder string) error {
	absPath, err := filepath.Abs(inputFolder)
	if err != nil {
		return err
	}

	if fi, err := os.Stat(absPath); os.IsNotExist(err) {
		return errors.New("Path " + absPath + " does not exists")
		if !fi.IsDir() {
			return errors.New(absPath + "is not folder")
		}
	}

	return nil
}
