package packer

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/itchio/lzma"
)

func TestSimplecompress(T *testing.T) {
	filename := "../../../test/output/simplecompress.dat"
	fd, _ := filepath.Abs(filename)
	defer os.Remove(fd)
	output, err := openOutput(filename)

	inputFolder := `../../../test/testdata/simplefolder`
	_, _, err = readInputFolder(inputFolder)

	T.Logf("Output struct: %v", output)
	written := closeOutput()

	T.Logf("Output file size: %d", written)
	if err != nil {
		T.Fatal(err)
	}

	T.Logf("Offsets table: %v", common.GetOffsets())

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		T.Fatal(err)
	}

	var b bytes.Buffer
	r := lzma.NewReader(f)
	io.CopyN(&b, r, int64(written))
	r.Close()
	readed := len(b.Bytes())
	if readed != 6 {
		T.Errorf("Unexpected uncompressed data size %d", readed)
	}

}
