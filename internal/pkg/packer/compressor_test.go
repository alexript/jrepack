package packer

import (
	//	"encoding/json"
	"bytes"
	"io"

	"os"
	"testing"

	common "github.com/alexript/jrepack/internal/pkg/common"
	"github.com/itchio/lzma"
)

func TestSimplecompress(T *testing.T) {
	filename := "../../../test/output/simplecompress.dat"
	output, err := openOutput(filename)
	defer os.Remove(output.File.Name())
	inputFolder := `../../../test/testdata/simplefolder`
	_, _, err = readInputFolder(inputFolder)
	closeOutput()
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
	io.Copy(&b, r)
	r.Close()
	if len(b.Bytes()) != 6 {
		T.Error("Unexpected uncompressed data size")
	}

}

//func TestHeader(T *testing.T) {
//	inputFolder := `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_stable`
//	dirinfo, rootfolder, _ := readInputFolder(inputFolder)
//	dump, _ := json.Marshal(dirinfo)
//	ioutil.WriteFile("../../../test/output/dirinfo.json", dump, 0644)
//	dump, _ = json.Marshal(rootfolder)
//	ioutil.WriteFile("../../../test/output/rootfolder.json", dump, 0644)

//}

//func TestFast(T *testing.T) {
//	inputFolder := `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_stable`
//	_, err := openOutput("../../../test/output/runtime.dat")
//	_, _, err = readInputFolder(inputFolder)
//	closeOutput()
//	if err != nil {
//		T.Fatal(err)
//	}
//}
