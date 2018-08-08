package packer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPack(T *testing.T) {
	filename := "../../../test/output/packtest2.dat"
	inputFolder := `../../../test/testdata/simplefolder`
	f, _ := filepath.Abs(filename)
	os.Remove(f)

	err := Pack("somefolder", "test", false)
	if err == nil {
		T.Fatal("Accept not existed input folder")
	}
	T.Log(err)

	err = Pack(".", "../../../test/output/.do_not_remove", false)
	if err == nil {
		T.Fatal("Accept existed output file")
	}
	T.Log(err)

	defer os.Remove(f)

	err = Pack(inputFolder, filename, false)
	if err != nil {
		T.Fatal(err)
	}

}

func TestPackZerosizeFile(T *testing.T) {
	filename := "../../../test/output/packtest3.dat"
	inputFolder := `../../../test/testdata/zerosize`
	f, _ := filepath.Abs(filename)
	os.Remove(f)

	defer os.Remove(f)

	err := Pack(inputFolder, filename, false)
	if err != nil {
		T.Fatal(err)
	}

}
