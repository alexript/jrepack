package unpacker

import (
	"testing"

	common "github.com/alexript/jrepack/internal/pkg/common"
)

func TestGetOutputPath(T *testing.T) {
	f1 := common.NewFolder("f1", false)
	f2 := common.NewFolder("f2", false)
	f3 := common.NewFolder("f3", false)
	f4 := common.NewFolder("f4", false)

	h := common.NewHeader(1000)
	f1id := h.Fold(0, &f1)
	f2id := h.Fold(f1id, &f2)
	f3id := h.Fold(f2id, &f3)
	f4id := h.Fold(f3id, &f4)

	T.Logf("Header: %v", h)

	dirpath, archpath, err := GetOutputPath(h, "", f4id)
	if err != nil {
		T.Error(err)
	}

	if archpath != nil {
		T.Errorf("Non-nil archive path %s", *archpath)
	}

	if *dirpath != "f1/f2/f3/f4" {
		T.Errorf("Wrong path %s", *dirpath)
	}

	dirpath, archpath, err = GetOutputPath(h, "", 0)
	if err != nil {
		T.Error(err)
	}

	if archpath != nil {
		T.Errorf("Non-nil archive path %s", *archpath)
	}

	if *dirpath != "" {
		T.Errorf("Wrong path %s", *dirpath)
	}

	dirpath, archpath, err = GetOutputPath(h, "", f1id)
	if err != nil {
		T.Error(err)
	}

	if archpath != nil {
		T.Errorf("Non-nil archive path %s", *archpath)
	}

	if *dirpath != "f1" {
		T.Errorf("Wrong path %s", *dirpath)
	}
}

func TestGetOutputPathWithArchive(T *testing.T) {
	f1 := common.NewFolder("f1", false)
	f2 := common.NewFolder("f2", false)
	f3 := common.NewFolder("f3", true)
	f4 := common.NewFolder("f4", false)

	h := common.NewHeader(1000)
	f1id := h.Fold(0, &f1)
	f2id := h.Fold(f1id, &f2)
	f3id := h.Fold(f2id, &f3)
	f4id := h.Fold(f3id, &f4)

	T.Logf("Header: %v", h)

	dirpath, archpath, err := GetOutputPath(h, "", f4id)
	if err != nil {
		T.Error(err)
	}

	if archpath == nil {
		T.Fatal("Nil archive path")
	}

	if *dirpath != "f1/f2/f3" || *archpath != "f4" {
		T.Errorf("Wrong path %s:%s", *dirpath, *archpath)
	}

}
