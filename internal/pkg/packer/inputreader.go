package packer

import (
	"archive/zip"
	"bytes"
	"errors"

	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	common "github.com/alexript/jrepack/internal/pkg/common"
)

func readInputFolder(inputFolder string) (*common.Dirinfo, *common.Folder, error) {
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
			common.AddFolderToFolder(parent, &subfolder)
			walkInputTree(fullname, &subfolder)
		} else {

			_, isContainer := common.IsContainer(fullname)

			if isContainer {
				subfolder := common.NewFolder(name, true)

				err := readContainer(&subfolder, fullname)
				if err != nil {
					return err
				}
				common.AddFolderToFolder(parent, &subfolder)

			} else {

				fileData, err := ioutil.ReadFile(fullname)
				if err != nil {
					return err
				}
				file, isNewHash := common.NewFile(name, fileData)
				common.AddFileToFolder(parent, file)
				if isNewHash {
					offset, _, err := compress(fileData)
					if err != nil {
						return err
					}
					common.SetOffset(uint32(offset), file.Hashsum)
				}
			}
		}
	}

	return nil
}

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
			common.AddFolderToFolder(container, &folder)

		} else {

			var b bytes.Buffer

			_, err = io.Copy(&b, rc)
			if err != nil {
				return err
			}
			file, isNewHash := common.NewFile(f.Name, b.Bytes())
			common.AddFileToFolder(container, file)
			if isNewHash {
				offset, _, err := compress(b.Bytes())
				if err != nil {
					return err
				}
				common.SetOffset(uint32(offset), file.Hashsum)
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
