package unpacker

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	common "github.com/alexript/jrepack/internal/pkg/common"
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

	return nil, nil
}
