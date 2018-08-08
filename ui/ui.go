package ui

import (
	"fmt"
)

// UI interface definition ----------------------------------------
const (
	EVT_UNPACK_DONE = 1
	EVT_PACK_DONE   = 2
)

type Hash struct {
	File      string
	Size      int
	Hash      []byte
	IsNewHash bool
}

type Folder struct {
	IsContainer bool
	Name        string
}

type Compressed struct {
	Len   int
	Total uint32
}

type JrepackUi interface {
	Error(message string) // Display some error
	Fatal(message string) // Display fatal error
	Info(message string)  // Display some information
	OnEnd(eventid int)    // Event on the packing/unpacking end
	Hashed(info Hash)
	NewFolder(info Folder)
	Compress(info Compressed)
	Unpack(readedFolders, foldersNum int)
}

// default UI implementation --------------------------------------

type defaultUi struct {
}

func (ui defaultUi) Error(message string) {
	panic(message)
}

func (ui defaultUi) Fatal(message string) {
	panic(message)
}

func (ui defaultUi) Info(message string) {
	fmt.Println(message)
}

func (ui defaultUi) OnEnd(eventid int) {
	fmt.Printf("Process %d complete.\n", eventid)
}

func (ui defaultUi) Hashed(info Hash) {

}

func (ui defaultUi) NewFolder(info Folder) {

}

func (ui defaultUi) Compress(info Compressed) {

}

func (ui defaultUi) Unpack(readedFolders, foldersNum int) {

}

// UI implementation setter and getter ----------------------------

var (
	someUi JrepackUi = defaultUi{}
)

func Set(ui JrepackUi) {
	someUi = ui
}

func Current() JrepackUi {
	return someUi
}
