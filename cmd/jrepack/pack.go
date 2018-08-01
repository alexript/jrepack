package main

import (
	"fmt"

	"github.com/alexript/jrepack"
)

var (
	inputFolder = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_stable`
	outputFile  = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_stable.jre`
)

func main() {
	err := jrepack.Pack(inputFolder, outputFile)
	if err != nil {
		fmt.Errorf("jre pack error: %v", err)
		return
	}

	fmt.Println("Packed.")
}
