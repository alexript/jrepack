package unpacker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexript/jrepack/internal/pkg/common"
	"github.com/alexript/jrepack/internal/pkg/packer"
)

const (
	filename_test      string = "../../../test/output/packtest.dat"
	inputFolder_test   string = `../../../test/testdata/mixedfolder`
	outputDirRoot_test string = "../../../test/output/unpacked"
	outputDirName_test string = "mixedfolder"
)

func prepareTestData() error {
	dropTestData()
	err := packer.Pack(inputFolder_test, filename_test, false)
	return err
}

func dropTestData() {
	f, _ := filepath.Abs(filename_test)
	os.Remove(f)

	root, _ := filepath.Abs(outputDirRoot_test)

	dirName := filepath.Join(root, outputDirName_test)

	common.RemoveDirReq(dirName)
}

func TestUnpacker(T *testing.T) {
	err := prepareTestData()
	if err != nil {
		T.Fatal(err)
	}
	defer dropTestData()

	err = UnPack("not_existed", outputDirRoot_test)
	if err == nil {
		T.Error("Accepted not existed archive")
	}

	err = UnPack(filename_test, outputDirRoot_test)
	if err == nil {
		T.Error("Accepted existed target folder")
	}

	root, _ := filepath.Abs(outputDirRoot_test)
	dirName := filepath.Join(root, outputDirName_test)
	err = UnPack(filename_test, dirName)
	if err != nil {
		T.Fatal(err)
	}
}
