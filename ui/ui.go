// Copyright (C) 2018  Alexander Malyshev

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

/*
Package ui is the interface for UI
*/
package ui

import (
	"fmt"
)

// UI interface definition ----------------------------------------
const (
	EvtUnpackDone = 1
	EvtPackDone   = 2
)

// Hash is the type for UI Hashed function
type Hash struct {
	File      string
	Size      int
	Hash      []byte
	IsNewHash bool
}

// Folder is the type for UI NewFolder function
type Folder struct {
	IsContainer bool
	Name        string
}

// Compressed is the type for UI Compress function
type Compressed struct {
	Len   int
	Total uint32
}

// JrepackUI is the main UI interface
type JrepackUI interface {
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

type defaultUI struct {
}

func (ui defaultUI) Error(message string) {
	panic(message)
}

func (ui defaultUI) Fatal(message string) {
	panic(message)
}

func (ui defaultUI) Info(message string) {
	fmt.Println(message)
}

func (ui defaultUI) OnEnd(eventid int) {
	fmt.Printf("Process %d complete.\n", eventid)
}

func (ui defaultUI) Hashed(info Hash) {

}

func (ui defaultUI) NewFolder(info Folder) {

}

func (ui defaultUI) Compress(info Compressed) {

}

func (ui defaultUI) Unpack(readedFolders, foldersNum int) {

}

// UI implementation setter and getter ----------------------------

var (
	someUI JrepackUI = defaultUI{}
)

// Set the current UI implementation
func Set(ui JrepackUI) {
	someUI = ui
}

// Current UI implementation
func Current() JrepackUI {
	return someUI
}
