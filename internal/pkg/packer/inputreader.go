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

package packer

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	common "github.com/alexript/jrepack/internal/pkg/common"
)

/*
readInputFolder is entry point of inputreader.
*/
func readInputFolder(inputFolder string) (*common.Dirinfo, *common.Folder, error) {
	runtime.GC()
	absPath, err := filepath.Abs(inputFolder)
	if err != nil {
		return nil, nil, err
	}

	fi, err := os.Stat(absPath)

	if os.IsNotExist(err) {
		return nil, nil, errors.New("Path " + absPath + " does not exists")

	}
	if !fi.IsDir() {
		return nil, nil, errors.New(absPath + " is not folder")
	}

	common.ClearDirinfo()
	rootfolder := common.NewFolder("_root_", false)
	err = walkInputTree(absPath, &rootfolder)
	if err != nil {
		return nil, nil, err
	}

	return common.GetDirinfo(), &rootfolder, nil
}

/*
walkInputTree is recursive walker
*/
func walkInputTree(dirname string, parent *common.Folder) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	if parent == nil {
		return errors.New("No parent folder pointer for " + dirname)
	}

	for _, fi := range files {
		name := fi.Name()
		fullname := filepath.Join(dirname, name)
		if fi.IsDir() {
			subfolder := common.NewFolder(name, false)
			err = common.AddFolderToFolder(parent, &subfolder)
			if err != nil {
				return err
			}
			err = walkInputTree(fullname, &subfolder)
			if err != nil {
				return err
			}
		} else {

			_, isContainer := common.IsContainer(fullname)

			if isContainer {
				subfolder := common.NewFolder(name, true)

				err = readContainer(&subfolder, fullname)
				if err != nil {
					return err
				}
				err = common.AddFolderToFolder(parent, &subfolder)
				if err != nil {
					return err
				}

			} else {

				fileData, err := ioutil.ReadFile(fullname)
				if err != nil {
					return err
				}
				file, isNewHash := common.NewFile(name, fileData)
				err = common.AddFileToFolder(parent, file)
				if err != nil {
					return err
				}
				if isNewHash {

					if len(fileData) > 0 {
						offset, _, err := compress(fileData)
						if err != nil {
							return err
						}
						common.SetOffset(uint32(offset), file.Hashsum)
					}

				}
			}
		}
	}

	return nil
}

/*
readContainer is recursive zip-file reader
*/
func readContainer(container *common.Folder, filename string) error {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		if f.FileInfo().IsDir() {
			folder := common.NewFolder(f.Name, false)
			err = common.AddFolderToFolder(container, &folder)
			if err != nil {
				return err
			}

		} else {

			fileData, err := ioutil.ReadAll(rc)

			if err != nil {
				return err
			}

			file, isNewHash := common.NewFile(f.Name, fileData)
			err = common.AddFileToFolder(container, file)
			if err != nil {
				return err
			}
			if isNewHash {

				if len(fileData) > 0 {
					offset, _, err := compress(fileData)
					if err != nil {
						return err
					}
					common.SetOffset(uint32(offset), file.Hashsum)
				}
			}

		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
