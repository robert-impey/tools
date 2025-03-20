package sdlib

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetSdFolder(t *testing.T) {
	containingDir := t.TempDir()
	testFileName := filepath.Join(containingDir, "test.txt")

	testFile, _ := os.Create(testFileName)
	defer testFile.Close()

	fmt.Fprintf(testFile, "test\n")

	var sdDir = filepath.Join(containingDir, ".stay-deleted")

	var fetchedDir, _ = GetSdFolder(testFileName)

	if fetchedDir != sdDir {
		t.Error(`GetSdFolder(containingDir) != sdDir`)
	}
}

func TestSetGetAction(t *testing.T) {
	dir := t.TempDir()

	tfns := [...]string{"test.txt", "file with spaces.txt", "file with [].txt"}

	for _, tfn := range tfns {
		testGetSetFile(dir, tfn, t)
	}
}

func TestSetGetActionFunnyCharsInDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "dir with [weird ( chars")

	os.Mkdir(dir, 0755)

	tfn := "test.txt"
	testGetSetFile(dir, tfn, t)
}

func testGetSetFile(dir string, tfn string, t *testing.T) {
	tfp := filepath.Join(dir, tfn)

	tf, _ := os.Create(tfp)
	defer tf.Close()

	action := Keep
	err := SetActionForFile(tfp, action)
	if err != nil {
		t.Error(err)
	}

	sdfp, err := GetSdFile(tfp)
	if err != nil {
		t.Error(err)
	}

	gotAction, err := GetActionForFile(sdfp, dir, os.Stderr)
	if err != nil {
		t.Error(err)
	}

	if gotAction.File != tfp {
		t.Error(fmt.Sprintf("gotAction.File is '%s', expecting '%s'!", gotAction.File, tfp))
	}

	if gotAction.Action != action {
		t.Error(fmt.Sprintf("gotAction.Action: %s!", getStringForAction(gotAction.Action)))
	}
}

func TestFindingSdFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "dir with [globbing] (chars)")

	os.Mkdir(dir, 0755)

	tfn := "test.txt"
	tfp := filepath.Join(dir, tfn)

	tf, _ := os.Create(tfp)
	defer tf.Close()

	action := Delete
	err := SetActionForFile(tfp, action)
	if err != nil {
		t.Error(err)
	}

	sdPath := filepath.Join(dir, SdFolderName)

	sdFiles, err := FindSdFiles(sdPath)
	if err != nil {
		t.Error(err)
	}

	foundFounds := len(sdFiles)
	if foundFounds != 1 {
		t.Error(fmt.Sprintf("got %d sd files, expecting 1!", foundFounds))
	}
}
