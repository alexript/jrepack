package packer

import (
	"testing"
)

func TestPack(T *testing.T) {
	err := Pack("test", "test", "test")
	if err != nil {
		T.Error(err)
	}
}
