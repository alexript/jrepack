package packer

import (
	"testing"
)

func TestReadNotexistedInputFolder(T *testing.T) {
	_, _, err := readInputFolder("./non_existed_folder_name")
	if err == nil {
		T.Error("Can read non existed input folder")
	}
}

func TestReadExistedInputFolder(T *testing.T) {
	_, _, err := readInputFolder("../packer")
	if err != nil {
		T.Error(err)
	}
}

func TestReadInvalidInputFolder(T *testing.T) {
	_, _, err := readInputFolder(`\\\something`) // this is invalid filename for windows
	if err == nil {
		T.Error("Can read invalid input folder")
	}
}

func TestReadInputNotFolder(T *testing.T) {
	_, _, err := readInputFolder("packer.go")
	if err == nil {
		T.Error("Can read file as input folder")
	}
}

func TestReadSimplefolder(T *testing.T) {

	inputFolder := `../../../test/testdata/simplefolder`
	dirinfo, _, err := readInputFolder(inputFolder)
	if err != nil || dirinfo == nil {
		T.Errorf("Unable to read test data from %v", inputFolder)
	}

	resultDirinfoLength := len(*dirinfo)
	if resultDirinfoLength != 1 {
		T.Errorf("Expected only 1 record in dirinfo. Result: %v", resultDirinfoLength)
	}
}

func TestReadSimplecontainer(T *testing.T) {

	inputFolder := `../../../test/testdata/simplecontainer`
	dirinfo, rootFodler, err := readInputFolder(inputFolder)
	if err != nil || dirinfo == nil {
		T.Errorf("Unable to read test data from %v", inputFolder)
	}

	T.Errorf("Root folder: %v", *rootFodler)

	resultDirinfoLength := len(*dirinfo)
	if resultDirinfoLength != 1 {
		T.Errorf("Expected only 1 record in dirinfo. Result: %v", resultDirinfoLength)
	}
}
