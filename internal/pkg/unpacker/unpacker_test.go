package unpacker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexript/jrepack/internal/pkg/common"
	"github.com/alexript/jrepack/internal/pkg/packer"
)

const (
	filename      string = "../../../test/output/packtest.dat"
	inputFolder   string = `../../../test/testdata/simplefolder`
	outputDirRoot string = "../../../test/output/unpacked"
	outputDirName string = "simplefolder"
)

func prepareTestData() error {

	err := packer.Pack(inputFolder, filename)
	return err
}

func dropTestData() {
	f, _ := filepath.Abs(filename)
	os.Remove(f)

	root, _ := filepath.Abs(outputDirRoot)

	dirName := filepath.Join(root, outputDirName)

	common.RemoveDirReq(dirName)
}

func TestUnpacker(T *testing.T) {
	err := prepareTestData()
	if err != nil {
		T.Fatal(err)
	}
	defer dropTestData()

	err = UnPack("not_existed", outputDirRoot)
	if err == nil {
		T.Error("Accepted not existed archive")
	}

	err = UnPack(filename, outputDirRoot)
	if err == nil {
		T.Error("Accepted existed tager folder")
	}

	root, _ := filepath.Abs(outputDirRoot)
	dirName := filepath.Join(root, outputDirName)
	err = UnPack(filename, dirName)
	if err != nil {
		T.Fatal(err)
	}
}
