package jrepack

import (
	"github.com/alexript/jrepack/internal/pkg/packer"
	"github.com/alexript/jrepack/internal/pkg/unpacker"
)

func Pack(inputFolder, outputFile string, dumpheader bool) error {
	return packer.Pack(inputFolder, outputFile, dumpheader)
}

func UnPack(inputFile, outputFolder string) error {
	return unpacker.UnPack(inputFile, outputFolder)
}
