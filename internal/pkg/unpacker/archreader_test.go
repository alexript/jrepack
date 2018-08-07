package unpacker

import (
	"testing"
)

func TestTest(T *testing.T) {
	err := prepareTestData()
	if err != nil {
		T.Fatal(err)
	}
	defer dropTestData()

	header, err := readArch(filename_test)
	if err != nil {
		T.Fatal(err)
	}
	T.Logf("Unpacked header: %v", header)

	if len(header.Folders) != 21 {
		T.Errorf("Unexpected number of folders: %v", len(header.Folders))
	}

}
