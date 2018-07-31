package main

import (
	"fmt"

	"github.com/alexript/jrepack"
)

var (
	outputFolder = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_unpacked`
	inputFile    = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172.jre`
)

func main() {

	err := jrepack.UnPack(inputFile, outputFolder)
	if err != nil {
		fmt.Errorf("jre unpack error: %v", err)
		return
	}

	fmt.Println("UnPacked.")

}
