package common

import (
	"encoding/hex"
	"testing"
)

var containerTypesTests = []struct {
	name        string
	expected    bool
	expectedExt string
}{
	{"some", false, ""},
	{"some.", false, ""},
	{"some.z", false, ""},
	{"some.zip", true, zipExt},
	{"some.Zip", true, zipExt},
	{"some.ZIp", true, zipExt},
	{"some.ZIP", true, zipExt},
	{"some.zIP", true, zipExt},
	{"some.ziP", true, zipExt},
	{"some.zIp", true, zipExt},
	{"some.piz", false, ""},
	{"some.jar", true, jarExt},
	{"some.Jar", true, jarExt},
	{"some.JAr", true, jarExt},
	{"some.JAR", true, jarExt},
	{"some.jAR", true, jarExt},
	{"some.jaR", true, jarExt},
	{"some.jAr", true, jarExt},
}

func TestIsContainer(t *testing.T) {
	for _, tt := range containerTypesTests {
		t.Run(tt.name, func(t *testing.T) {
			ct, result := isContainer(tt.name)
			if tt.expected != result {
				t.Errorf("Result: %v, Expected: %v", result, tt.expected)
			}
			if result && ct == nil {
				t.Error("Result is true, but container is nil")
			}
			if !result && ct != nil {
				t.Error("Result is false, but container is not nil")
			}
			if result && ct.Extension != tt.expectedExt {
				t.Errorf("Result extension: %v, expected: %v", ct.Extension, tt.expectedExt)
			}
		})
	}
}

func fromHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func toHex(b []byte) string {
	return hex.EncodeToString(b)
}

func TestNewFile(t *testing.T) {
	expectedName := "test"
	expectedHexString := "010203040506"
	expectedBody := fromHex(expectedHexString)
	expectedLength := len(expectedBody)
	expectedHash := "55b88037ec60704aa5dc318200f6998afb24b31d1a7b1d3f5d2263472ea73f70"
	f := newFile(expectedName, expectedBody)
	if expectedName != f.Name {
		t.Errorf("Result: '%v', expected: '%v'", f.Name, expectedName)
	}
	if expectedLength != f.Size {
		t.Errorf("Result: '%v', expected: '%v'", f.Size, expectedLength)
	}
	resultHash := toHex(f.Hashsum)
	if expectedHash != resultHash {
		t.Errorf("Result: '%v', expected: '%v'", resultHash, expectedHash)
	}
}

func TestNewFolder(t *testing.T) {
	expectedName := "test"
	expectedZeroNum := 0
	f := newFolder(expectedName)
	if f.Name != expectedName {
		t.Errorf("Result: '%v', expected: '%v'", f.Name, expectedName)
	}
	resultFilesNum := len(f.Files)
	if expectedZeroNum != resultFilesNum {
		t.Errorf("Result: '%v', expected: '%v'", resultFilesNum, expectedZeroNum)
	}
	resultFoldersNum := len(f.Folders)
	if expectedZeroNum != resultFoldersNum {
		t.Errorf("Result: '%v', expected: '%v'", resultFoldersNum, expectedZeroNum)
	}
	resultContainersNum := len(f.Containers)
	if expectedZeroNum != resultContainersNum {
		t.Errorf("Result: '%v', expected: '%v'", resultContainersNum, expectedZeroNum)
	}
}

func TestNewContainer(t *testing.T) {
	c := newContainer("test")
	if c != nil {
		t.Error("Accepted not container name")
	}

	c = newContainer("test.zip")
	if c == nil {
		t.Error("ZIP container is not accepted")
	}

	expectedName := "test.jar"
	c = newContainer(expectedName)
	if c == nil {
		t.Error("JAR container is not accepted")
	}

	if expectedName != c.Content.Name {
		t.Error("Container have no name")
	}
}

func TestAddFileToFolder(t *testing.T) {
	expectedHexString := "010203040506"
	expectedBody := fromHex(expectedHexString)
	err := addFileToFolder(nil, nil)
	if err == nil {
		t.Error("Nil folder and Nil file are accepted")
	}

	f := newFile("test", expectedBody)
	err = addFileToFolder(nil, &f)
	if err == nil {
		t.Error("Nil folder are accepted")
	}

	fold := newFolder("test")
	err = addFileToFolder(&fold, nil)
	if err == nil {
		t.Error("Nil file are accepted")
	}

	err = addFileToFolder(&fold, &f)
	if err != nil {
		t.Error("Unable to append file to folder")
	}

	expectedSize := 1
	resultSize := len(fold.Files)
	if expectedSize != resultSize {
		t.Errorf("Expected %v files in folder. Result: %v", expectedSize, resultSize)
	}
}

func TestAddFolderToFolder(t *testing.T) {

	err := addFolderToFolder(nil, nil)
	if err == nil {
		t.Error("Nil folder and Nil folder are accepted")
	}

	src := newFolder("test")
	err = addFolderToFolder(nil, &src)
	if err == nil {
		t.Error("Nil destination folder are accepted")
	}

	dest := newFolder("test")
	err = addFolderToFolder(&dest, nil)
	if err == nil {
		t.Error("Nil source folder are accepted")
	}

	err = addFolderToFolder(&dest, &src)
	if err != nil {
		t.Error("Unable to append fodler to folder")
	}

	expectedSize := 1
	resultSize := len(dest.Folders)
	if expectedSize != resultSize {
		t.Errorf("Expected %v fodlers in folder. Result: %v", expectedSize, resultSize)
	}
}

func TestAddContainerToFolder(t *testing.T) {

	err := addContainerToFolder(nil, nil)
	if err == nil {
		t.Error("Nil folder and Nil container are accepted")
	}

	dest := newFolder("test")
	err = addContainerToFolder(&dest, nil)
	if err == nil {
		t.Error("Nil container are accepted")
	}

	c := newContainer("test.zip")
	err = addContainerToFolder(nil, c)
	if err == nil {
		t.Error("Nil destination folder are accepted")
	}

	err = addContainerToFolder(&dest, c)
	if err != nil {
		t.Error("Unable to append container to folder")
	}

	expectedSize := 1
	resultSize := len(dest.Containers)
	if expectedSize != resultSize {
		t.Errorf("Expected %v containers in folder. Result: %v", expectedSize, resultSize)
	}
}

func TestAddFileToContainer(t *testing.T) {
	expectedHexString := "010203040506"
	expectedBody := fromHex(expectedHexString)
	err := addFileToContainer(nil, nil)
	if err == nil {
		t.Error("Nil container and Nil file are accepted")
	}

	f := newFile("test", expectedBody)
	err = addFileToContainer(nil, &f)
	if err == nil {
		t.Error("Nil container are accepted")
	}

	c := newContainer("test.zip")
	err = addFileToContainer(c, nil)
	if err == nil {
		t.Error("Nil file are accepted")
	}

	err = addFileToContainer(c, &f)
	if err != nil {
		t.Error("Unable to append file to container")
	}

	expectedSize := 1
	resultSize := len(c.Content.Files)
	if expectedSize != resultSize {
		t.Errorf("Expected %v files in container. Result: %v", expectedSize, resultSize)
	}
}

func TestAddFolderToContainer(t *testing.T) {

	err := addFolderToContainer(nil, nil)
	if err == nil {
		t.Error("Nil container and Nil folder are accepted")
	}

	src := newFolder("test")
	err = addFolderToContainer(nil, &src)
	if err == nil {
		t.Error("Nil contaienr are accepted")
	}

	c := newContainer("test.zip")
	err = addFolderToContainer(c, nil)
	if err == nil {
		t.Error("Nil folder are accepted")
	}

	err = addFolderToContainer(c, &src)
	if err != nil {
		t.Error("Unable to append fodler to container")
	}

	expectedSize := 1
	resultSize := len(c.Content.Folders)
	if expectedSize != resultSize {
		t.Errorf("Expected %v fodlers in container. Result: %v", expectedSize, resultSize)
	}
}

func TestAddContainerToContainer(t *testing.T) {

	err := addContainerToContainer(nil, nil)
	if err == nil {
		t.Error("Nil container and Nil container are accepted")
	}

	dest := newContainer("test.zip")
	err = addContainerToContainer(dest, nil)
	if err == nil {
		t.Error("Nil source container are accepted")
	}

	c := newContainer("test.zip")
	err = addContainerToContainer(nil, c)
	if err == nil {
		t.Error("Nil destination container are accepted")
	}

	err = addContainerToContainer(dest, c)
	if err != nil {
		t.Error("Unable to append container to container")
	}

	expectedSize := 1
	resultSize := len(dest.Content.Containers)
	if expectedSize != resultSize {
		t.Errorf("Expected %v containers in container. Result: %v", expectedSize, resultSize)
	}
}

func TestAddFileToDirinfo(t *testing.T) {

	// some file#1
	expectedName1 := "test1"
	expectedHexString1 := "010203040506"
	expectedBody1 := fromHex(expectedHexString1)

	// file#2 have same content as file#1 but different name
	expectedName2 := "test2"
	expectedHexString2 := "010203040506"
	expectedBody2 := fromHex(expectedHexString2)

	// file#3 have the same name as file#1 but different content
	expectedName3 := "test1"
	expectedHexString3 := "01020304050607"
	expectedBody3 := fromHex(expectedHexString3)

	// file#4 is unique
	expectedName4 := "test3"
	expectedHexString4 := "0102030405060708"
	expectedBody4 := fromHex(expectedHexString4)

	// dummy folder to collect files
	fold := newFolder("dummy")

	clearDirinfo()
	currentLen := len(dirinfo)
	if currentLen > 0 {
		t.Error("Unable to cleanup the dirinfo")
	}

	f1 := newFile(expectedName1, expectedBody1)
	addFileToFolder(&fold, &f1)
	currentLen = len(dirinfo)
	if currentLen != 1 {
		t.Error("Unable to add file into empty dirinfo")
	}

	f2 := newFile(expectedName1, expectedBody1)
	addFileToFolder(&fold, &f2)
	currentLen = len(dirinfo)
	if currentLen != 1 {
		t.Error("Different key is produced for the same file")
	}

	f3 := newFile(expectedName2, expectedBody2)
	addFileToFolder(&fold, &f3)
	currentLen = len(dirinfo)
	if currentLen != 1 {
		t.Error("Same file content but different name failed")
	}

	f4 := newFile(expectedName3, expectedBody3)
	addFileToFolder(&fold, &f4)
	currentLen = len(dirinfo)
	if currentLen != 2 {
		t.Error("Different content not separated")
	}

	f5 := newFile(expectedName4, expectedBody4)
	addFileToFolder(&fold, &f5)
	currentLen = len(dirinfo)
	if currentLen != 3 {
		t.Error("Totally different file not separated")
	}
}
