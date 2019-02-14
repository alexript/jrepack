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

package cmdui

import (
	"fmt"

	"github.com/alexript/jrepack/ui"
)

// CommandlineUI struct is concrete and simple UI implementation
type CommandlineUI struct {
	Archivefile string
}

// Error will produce panic
func (u CommandlineUI) Error(message string) {
	panic(message)
}

// Fatal will produce panic
func (u CommandlineUI) Fatal(message string) {
	panic(message)
}

// Info will Println message
func (u CommandlineUI) Info(message string) {
	fmt.Println(message)
}

// OnEnd will Println simple info message
func (u CommandlineUI) OnEnd(eventid int) {
	switch eventid {
	case ui.EvtUnpackDone:
		fmt.Println("Unpacking complete.")

	case ui.EvtPackDone:
		fmt.Println("Packing complete.")

	default:
		fmt.Println("Process complete.")

	}
}

// Hashed will do nothing
func (u CommandlineUI) Hashed(info ui.Hash) {

}

// NewFolder will do nothing
func (u CommandlineUI) NewFolder(info ui.Folder) {

}

// Compress will do nothing
func (u CommandlineUI) Compress(info ui.Compressed) {

}

var percentage int

// Unpack will Print compressed files percentage on percents change
func (u CommandlineUI) Unpack(readedFolders, foldersNum int) {
	p := int(readedFolders * 100 / foldersNum)
	if p != percentage {
		fmt.Printf("Unpacking: %d%%\n", p)
		percentage = p
	}
}
