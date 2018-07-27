package packer

import (
	"testing"
)

func TestReadNotexistedInputFolder(T *testing.T) {
	err := readInputFolder("./non_existed_folder_name")
	if err == nil {
		T.Error("Can read non existed input folder")
	}
}

func TestReadExistedInputFolder(T *testing.T) {
	err := readInputFolder("../packer")
	if err != nil {
		T.Error(err)
	}
}

func TestReadInvalidInputFolder(T *testing.T) {
	err := readInputFolder(`\\\something`) // this is invalid filename for windows
	if err == nil {
		T.Error("Can read invalid input folder")
	}
}

func TestReadInputNotFolder(T *testing.T) {
	err := readInputFolder("packer.go")
	if err != nil {
		T.Error(err)
	}
}
