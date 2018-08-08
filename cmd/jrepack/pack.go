package main

import (
	"fmt"

	"github.com/alexript/jrepack"
	"github.com/alexript/jrepack/ui"
)

var (
	inputFolder = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_stable`
	outputFile  = `D:\workspace\ETax-2.0\runtimes\runtime\jre8u172_stable.jre`
)

func main() {
	err := jrepack.Pack(inputFolder, outputFile, false)
	if err != nil {
		ui.Current().Error(fmt.Sprintf("jre pack error: %v", err))

		return
	}
	ui.Current().OnEnd("Packed.")

}
