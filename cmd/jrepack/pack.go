package main

import (
	"fmt"

	"github.com/alexript/jrepack"
)

var (
	inputFolder = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172`
	outputFile  = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172.jre`
	workDir     = `D:\workspace\ETax-2.0\runtimes\runtime`
)

func main() {
	err := jrepack.Pack(inputFolder, outputFile, workDir)
	if err != nil {
		fmt.Errorf("jre pack error: %v", err)
		return
	}

	fmt.Println("Packed.")
}