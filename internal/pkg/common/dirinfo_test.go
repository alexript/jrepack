package common

import (
	"encoding/hex"
	"encoding/json"
	"path"
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
			ct, result := IsContainer(tt.name)
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
	f, _ := NewFile(expectedName, expectedBody)
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
	f := NewFolder(expectedName, false)
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

}

func TestAddFileToFolder(t *testing.T) {
	expectedHexString := "010203040506"
	expectedBody := fromHex(expectedHexString)
	err := AddFileToFolder(nil, nil)
	if err == nil {
		t.Error("Nil folder and Nil file are accepted")
	}

	f, _ := NewFile("test", expectedBody)
	err = AddFileToFolder(nil, f)
	if err == nil {
		t.Error("Nil folder are accepted")
	}

	fold := NewFolder("test", false)
	err = AddFileToFolder(&fold, nil)
	if err == nil {
		t.Error("Nil file are accepted")
	}

	err = AddFileToFolder(&fold, f)
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

	err := AddFolderToFolder(nil, nil)
	if err == nil {
		t.Error("Nil folder and Nil folder are accepted")
	}

	src := NewFolder("test", false)
	err = AddFolderToFolder(nil, &src)
	if err == nil {
		t.Error("Nil destination folder are accepted")
	}

	dest := NewFolder("test", false)
	err = AddFolderToFolder(&dest, nil)
	if err == nil {
		t.Error("Nil source folder are accepted")
	}

	err = AddFolderToFolder(&dest, &src)
	if err != nil {
		t.Error("Unable to append fodler to folder")
	}

	expectedSize := 1
	resultSize := len(dest.Folders)
	if expectedSize != resultSize {
		t.Errorf("Expected %v fodlers in folder. Result: %v", expectedSize, resultSize)
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
	fold := NewFolder("dummy", false)

	ClearDirinfo()
	currentLen := len(dirinfo)
	if currentLen > 0 {
		t.Error("Unable to cleanup the dirinfo")
	}

	f1, isNewHash1 := NewFile(expectedName1, expectedBody1)
	AddFileToFolder(&fold, f1)
	currentLen = len(dirinfo)
	if currentLen != 1 || !isNewHash1 {
		t.Error("Unable to add file into empty dirinfo")
	}

	f2, isNewHash2 := NewFile(expectedName1, expectedBody1)
	AddFileToFolder(&fold, f2)
	currentLen = len(dirinfo)
	if currentLen != 1 || isNewHash2 {
		t.Error("Different key is produced for the same file")
	}

	f3, isNewHash3 := NewFile(expectedName2, expectedBody2)
	AddFileToFolder(&fold, f3)
	currentLen = len(dirinfo)
	if currentLen != 1 || isNewHash3 {
		t.Error("Same file content but different name failed")
	}

	f4, isNewHash4 := NewFile(expectedName3, expectedBody3)
	AddFileToFolder(&fold, f4)
	currentLen = len(dirinfo)
	if currentLen != 2 || !isNewHash4 {
		t.Error("Different content not separated")
	}

	f5, isNewHash5 := NewFile(expectedName4, expectedBody4)
	AddFileToFolder(&fold, f5)
	currentLen = len(dirinfo)
	if currentLen != 3 || !isNewHash5 {
		t.Error("Totally different file not separated")
	}
}

func TestHasFolder(T *testing.T) {
	parent := NewFolder("parent", false)
	child := NewFolder("child", false)
	AddFolderToFolder(&parent, &child)

	f, err := (&parent).HasFolder("child")
	if err != nil {
		T.Error(err)
	}
	if f == nil {
		T.Fatal("Unable to find existing folder")
	}
	if f.Name != "child" {
		T.Error("Unexpected child name " + f.Name)
	}

	f, err = (&parent).HasFolder("nonexisted")
	if err != nil {
		T.Error(err)
	}
	if f != nil {
		T.Error("found not existed child")
	}

}

func TestMkdirAll(T *testing.T) {
	path := "/f1/f2"
	root := NewFolder("root", false)
	result, err := MkdirAll(&root, path)
	if err != nil {
		T.Fatal(err)
	}
	if result.Name != "f2" {
		T.Error("Result is not f2 but " + result.Name)
	}

	f1level := root.Folders
	if len(f1level) != 1 {
		T.Fatal("Not exacly one folder on f1 level")
	}

	f1 := f1level[0]
	f1name := f1.Name
	if f1name != "f1" {
		T.Error("f1 folder is " + f1name)
	}

	f2level := f1.Folders
	if len(f2level) != 1 {
		T.Fatal("Not exactly one folder on f2 level")
	}
	f2 := f2level[0]
	f2name := f2.Name
	if f2name != "f2" {
		T.Error("f2 folder have wrong name " + f2name)
	}

	// repeate

	result, err = MkdirAll(&root, path)
	if err != nil {
		T.Fatal(err)
	}

	f1level = root.Folders
	if len(f1level) != 1 {
		T.Fatal("Not exacly one folder on f1 level")
	}

	f1 = f1level[0]

	f2level = f1.Folders
	if len(f2level) != 1 {
		T.Fatal("Not exactly one folder on f2 level")
	}
	dump, _ := json.MarshalIndent(root, "// ", "   ")
	T.Logf("Root folder: %s", dump)
}

func TestAddFoldersAndFilesToFolder(T *testing.T) {
	testCase := "f1/f2/file.txt"
	hexString := "010203040506"
	body := fromHex(hexString)
	parent := NewFolder("test", false)
	file, _ := NewFile(testCase, body)
	AddFileToFolder(&parent, file)

	level1 := parent.Folders
	if len(level1) != 1 {
		T.Fatal("f1 folder is not created")
	}

	f1 := level1[0]
	level2 := f1.Folders
	if len(level2) != 1 {
		T.Fatal("f2 folder is not created")
	}

	f2 := level2[0]
	files := f2.Files
	if len(files) != 1 {
		T.Fatal("file.txt is not created")
	}

	f := files[0]
	if f.Name != "file.txt" {
		T.Error("file name is not trimmed")
	}

	dump, _ := json.MarshalIndent(parent, "// ", "   ")
	T.Logf("Root folder: %s", dump)
}

func TestAddContainerToFolder(T *testing.T) {
	// simulate zip structure: dirs has '/' at end
	expectedHexString := "010203040506"
	expectedBody := fromHex(expectedHexString)

	s1 := "/F/"
	s2 := "/F/test"
	zip := NewFolder("zip", true)
	folder := NewFolder(s1, false)
	file, _ := NewFile(s2, expectedBody)
	AddFolderToFolder(&zip, &folder)
	AddFileToFolder(&zip, file)

	foldersInF := zip.Folders[0].Folders

	log := T.Logf
	if len(foldersInF) != 0 {
		log := T.Errorf
		log("Folder F contain subfolders")
	}

	dirname := path.Dir(s1)
	basename := path.Base(s1)
	log("Dirname: %s, Basename: %s", dirname, basename)
	log("Struct: %v", zip)
}
