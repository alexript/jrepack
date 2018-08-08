package cmdui

import (
	"fmt"

	"github.com/alexript/jrepack/ui"
)

type CommandlineUi struct {
	Archivefile string
}

func (u CommandlineUi) Error(message string) {
	panic(message)
}

func (u CommandlineUi) Fatal(message string) {
	panic(message)
}

func (u CommandlineUi) Info(message string) {
	fmt.Println(message)
}

func (u CommandlineUi) OnEnd(eventid int) {
	switch eventid {
	case ui.EVT_UNPACK_DONE:
		fmt.Println("Unpacking complete.")

	case ui.EVT_PACK_DONE:
		fmt.Println("Packing complete.")

	default:
		fmt.Println("Process complete.")

	}
}

func (u CommandlineUi) Hashed(info ui.Hash) {

}

func (u CommandlineUi) NewFolder(info ui.Folder) {

}

func (u CommandlineUi) Compress(info ui.Compressed) {

}

var percentage int = 0

func (u CommandlineUi) Unpack(readedFolders, foldersNum int) {
	p := int(readedFolders * 100 / foldersNum)
	if p != percentage {
		fmt.Printf("Unpacking: %d%%\n", p)
		percentage = p
	}
}
