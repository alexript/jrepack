package main

import (
	"fmt"

	"github.com/alexript/jrepack"
	"github.com/alexript/jrepack/ui"
)

var (
	outputFolder = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_unpacked`
	inputFile    = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_stable.jre`
)

func main() {

	err := jrepack.UnPack(inputFile, outputFolder)
	if err != nil {
		ui.Current().Error(fmt.Sprintf("jre unpack error: %v", err))
		return
	}
	ui.Current().OnEnd("UnPacked.")

}
