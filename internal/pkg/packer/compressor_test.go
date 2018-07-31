package packer

import (
	"os"
	"testing"
)

func TestSimplecompress(T *testing.T) {

	output, err := openOutput("../../../test/output/simplecompress.dat")
	defer os.Remove(output.File.Name())
	inputFolder := `../../../test/testdata/simplefolder`
	_, _, err = readInputFolder(inputFolder)
	closeOutput()
	if err != nil {
		T.Fatal(err)
	}
	fi, _ := os.Stat(output.File.Name())
	if fi.Size() != 6 {
		T.Errorf("Unexpected size of the output: %d", fi.Size())
	}
}
