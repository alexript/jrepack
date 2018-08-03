package packer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPack(T *testing.T) {
	filename := "../../../test/output/packtest.dat"
	inputFolder := `../../../test/testdata/simplefolder`
	f, _ := filepath.Abs(filename)
	os.Remove(f)

	err := Pack("somefolder", "test")
	if err == nil {
		T.Fatal("Accept not existed input folder")
	}
	T.Log(err)

	err = Pack(".", "../../../test/output/.do_not_remove")
	if err == nil {
		T.Fatal("Accept existed output file")
	}
	T.Log(err)

	defer os.Remove(f)

	err = Pack(inputFolder, filename)
	if err != nil {
		T.Fatal(err)
	}

}
