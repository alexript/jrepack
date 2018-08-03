package unpacker

import (
	"bytes"
	"io"
	"os"
	"runtime"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/itchio/lzma"
)

func Decompress(header *common.Header, filename string, output string) error {

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		T.Fatal(err)
	}

	needToRead := header.Size

	var b bytes.Buffer
	r := lzma.NewReader(f)

	for _, dataRecord := range header.Data {
		b.Reset()
		io.CopyN(&b, r, int64(dataRecord.Size))
	}
	r.Close()
	b.Reset()

	defer runtime.GC()
	readed := len(b.Bytes())
	if readed != 12 {
		T.Errorf("Unexpected uncompressed data size %d", readed)
	}

	return nil
}
