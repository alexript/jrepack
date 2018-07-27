package jrepack

import (
	"github.com/alexript/jrepack/internal/pkg/packer"
	"github.com/alexript/jrepack/internal/pkg/unpacker"
)

func Pack(inputFolder, outputFile, workFolder string) error {
	return packer.Pack(inputFolder, outputFile, workFolder)
}

func UnPack(inputFile, outputFolder, workFolder string) error {
	return unpacker.UnPack(inputFile, outputFolder, workFolder)
}
