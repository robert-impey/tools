package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

	gotAction, err := GetActionForFile(sdfp, dir)
	if err != nil {
		t.Error(err)
	}

	if gotAction.File != tfp {
		t.Errorf("gotAction.File is '%s', expecting '%s'!", gotAction.File, tfp)
	}

	if gotAction.Action != action {
		t.Errorf("gotAction.Action: %s!", getStringForAction(gotAction.Action))
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
		t.Errorf("got %d sd files, expecting 1!", foundFounds)
	}
}

func TestSweepFromDirectories(t *testing.T) {
	// Create two separate directories to sweep from
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	// In dir1, create a file we intend to delete
	delFile1 := filepath.Join(dir1, "delete_me.txt")
	if err := os.WriteFile(delFile1, []byte("x"), 0644); err != nil {
		t.Fatalf("write delFile1: %v", err)
	}
	if err := SetActionForFile(delFile1, Delete); err != nil {
		t.Fatalf("SetActionForFile(delFile1, Delete): %v", err)
	}

	// In dir2, create a file we intend to keep
	keepFile2 := filepath.Join(dir2, "keep_me.txt")
	if err := os.WriteFile(keepFile2, []byte("y"), 0644); err != nil {
		t.Fatalf("write keepFile2: %v", err)
	}
	if err := SetActionForFile(keepFile2, Keep); err != nil {
		t.Fatalf("SetActionForFile(keepFile2, Keep): %v", err)
	}

	// Call SweepFromDirectories on both dirs. expiryMonths only affects removal of old SD metadata,
	// not whether a file marked Delete is removed. we can use any reasonable value (e.g., 6).
	if err := SweepFromDirectories([]string{dir1, dir2}, 6, false); err != nil {
		t.Fatalf("SweepFromDirectories returned error: %v", err)
	}

	// Assert: delete_me.txt is gone
	if _, err := os.Stat(delFile1); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be deleted, got err=%v", delFile1, err)
	}

	// Assert: keep_me.txt still exists
	if _, err := os.Stat(keepFile2); err != nil {
		t.Fatalf("expected %s to be kept, stat err=%v", keepFile2, err)
	}
}

// Desired behavior: SweepFromDirectories should continue when a directory in the list is missing
// and process the remaining directories instead of stopping at the first error.
// It should also return aggregated errors and must not fail the whole run.
func TestSweepFromDirectories_ContinuesOnMissingDirectory(t *testing.T) {
	// dir1 and dir3 exist; dirMissing does not
	dir1 := t.TempDir()
	dir3 := t.TempDir()
	dirMissing := filepath.Join(t.TempDir(), "does-not-exist") // do not create

	// Mark a file in dir1 for deletion
	delFile1 := filepath.Join(dir1, "delete_me.txt")
	if err := os.WriteFile(delFile1, []byte("x"), 0644); err != nil {
		t.Fatalf("write delFile1: %v", err)
	}
	if err := SetActionForFile(delFile1, Delete); err != nil {
		t.Fatalf("SetActionForFile(delFile1, Delete): %v", err)
	}

	// Mark a file in dir3 for deletion (should still be processed even if one dir is missing)
	delFile3 := filepath.Join(dir3, "delete_me_too.txt")
	if err := os.WriteFile(delFile3, []byte("z"), 0644); err != nil {
		t.Fatalf("write delFile3: %v", err)
	}
	if err := SetActionForFile(delFile3, Delete); err != nil {
		t.Fatalf("SetActionForFile(delFile3, Delete): %v", err)
	}

	// Call SweepFromDirectories with a missing directory in the middle
	if err := SweepFromDirectories([]string{dir1, dirMissing, dir3}, 6, false); err != nil {
		// We expect an error corresponding to the missing directory.
		log.Printf("Expected error from SweepFromDirectories: %v\n", err)
	} else {
		t.Fatal("Expected SweepFromDirectories to return an error for the missing directory, but it returned nil")
	}

	// dir1's delete should have happened
	if _, err := os.Stat(delFile1); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be deleted, got err=%v", delFile1, err)
	}

	// dir3 should have been processed even though dirMissing did not exist
	if _, err := os.Stat(delFile3); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be deleted even with a missing directory present, got err=%v", delFile3, err)
	}
}

func TestSweepFromDirectories_NFCDeletionKeepsNFD(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("Unicode normalization filename test is not reliable on macOS")
	}

	dir := t.TempDir()

	// NFC: café.txt with precomposed U+00E9
	nfcName := "café.txt"

	// NFD: café.txt with 'e' + combining acute accent U+0301
	nfdName := "café.txt"

	nfcFile := filepath.Join(dir, nfcName)
	nfdFile := filepath.Join(dir, nfdName)

	if err := os.WriteFile(nfcFile, []byte("nfc"), 0644); err != nil {
		t.Fatalf("write nfcFile: %v", err)
	}
	if err := os.WriteFile(nfdFile, []byte("nfd"), 0644); err != nil {
		t.Fatalf("write nfdFile: %v", err)
	}

	if err := SetActionForFile(nfcFile, Delete); err != nil {
		t.Fatalf("SetActionForFile(nfcFile, Delete): %v", err)
	}
	if err := SetActionForFile(nfdFile, Keep); err != nil {
		t.Fatalf("SetActionForFile(nfdFile, Keep): %v", err)
	}

	if err := SweepFromDirectories([]string{dir}, 6, false); err != nil {
		t.Fatalf("SweepFromDirectories returned error: %v", err)
	}

	if _, err := os.Stat(nfcFile); !os.IsNotExist(err) {
		t.Fatalf("expected NFC file %q to be deleted, got err=%v", nfcFile, err)
	}

	if _, err := os.Stat(nfdFile); err != nil {
		t.Fatalf("expected NFD file %q to remain, stat err=%v", nfdFile, err)
	}
}
