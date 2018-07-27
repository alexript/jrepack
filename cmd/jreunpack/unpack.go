package main

import (
	"fmt"

	"github.com/alexript/jrepack"
)

var (
	outputFolder = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_unpacked`
	inputFile    = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172.jre`
	workDir      = `D:\workspace\ETax-2.0\runtimes\runtime`
)

func main() {

	err := jrepack.UnPack(inputFile, outputFolder, workDir)
	if err != nil {
		fmt.Errorf("jre unpack error: %v", err)
		return
	}

	fmt.Println("UnPacked.")

}
