package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDirsAndFiles(t *testing.T) {
	// Create a temporary directory structure for testing
	// tempDir/dir1/file1.txt
	// tempDir/dir1/file2.txt
	// tempDir/dir2/file3.txt
	// tempDir/file4.txt

	tempDir, err := os.MkdirTemp("", "test_tidy_folder")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir2")

	err = os.Mkdir(dir1, 0o755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}
	err = os.Mkdir(dir2, 0o755)
	if err != nil {
		t.Fatalf("Failed to create dir2: %v", err)
	}

	createFile(t, filepath.Join(dir1, "file1.txt"), "content1")
	createFile(t, filepath.Join(dir1, "file2.txt"), "content2")
	createFile(t, filepath.Join(dir2, "file3.txt"), "content3")
	createFile(t, filepath.Join(tempDir, "file4.txt"), "content4")

	// Call the function under test
	dirsAndFiles, err := BuildDirsAndFiles(tempDir)
	if err != nil {
		t.Fatalf("BuildDirsAndFiles failed: %v", err)
	}

	// Assertions
	if len(dirsAndFiles) != 4 { // tempDir, dir1, dir2 + the root directory itself if it's walked as a file (though unlikely here)
		t.Errorf("Expected at least 3 entries (files/dirs), got %d", len(dirsAndFiles))
	}

	// Check for presence of expected files by name, regardless of parent grouping
	foundFile4 := false
	foundFile1 := false
	foundFile2 := false
	foundFile3 := false

	for _, file := range dirsAndFiles {
		switch file.Name {
		case "file4.txt":
			if !foundFile4 {
				foundFile4 = true
			}
		case "file1.txt":
			if !foundFile1 {
				foundFile1 = true
			}
		case "file2.txt":
			if !foundFile2 {
				foundFile2 = true
			}
		case "file3.txt":
			if !foundFile3 {
				foundFile3 = true
			}
		}
	}

	if !foundFile4 {
		t.Errorf("Expected file 'file4.txt' not found in any entry")
	}
	if !foundFile1 {
		t.Errorf("Expected file 'file1.txt' not found in any entry")
	}
	if !foundFile2 {
		t.Errorf("Expected file 'file2.txt' not found in any entry")
	}
	if !foundFile3 {
		t.Errorf("Expected file 'file3.txt' not found in any entry")
	}
}

func createFile(t *testing.T, path, content string) {
	err := os.WriteFile(path, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", path, err)
	}
}

func TestFindMatchingStems(t *testing.T) {
	// Create a list of entries simulating the output of BuildDirsAndFiles
	dirsAndFiles := []DirEntry{
		{Path: "/tmp/test/a.txt", Name: "a.txt"},
		{Path: "/tmp/test/a1.txt", Name: "a1.txt"},
		{Path: "/tmp/test/b.txt", Name: "b.txt"},
		{Path: "/tmp/test/a.md", Name: "a.md"},
	}

	matches := FindMatchingStems(dirsAndFiles)

	if len(matches) < 1 {
		t.Fatalf("expected at least one match, got %d", len(matches))
	}

	// Check for the specific expected pair (order might vary due to iteration)
	foundMatch := false
	for _, pair := range matches {
		// We check if this pair represents a valid stem match: e.g., "a.txt" and "a1.txt"
		stem0, ext0, ok0 := splitStemExt(pair[0].Path)
		stem1, ext1, ok1 := splitStemExt(pair[1].Path)

		if ok0 && ok1 && ext0 == ext1 && strings.HasPrefix(stem1, stem0) && stem1 != stem0 {
			foundMatch = true
			break
		}
	}

	if !foundMatch {
		t.Errorf("Did not find expected matching stems pair (e.g., a.txt and a1.txt)")
	}
}

func TestReadDirectories(t *testing.T) {
	tempDir := t.TempDir()
	input := filepath.Join(tempDir, "dirs.txt")

	content := "  /alpha  \n\n# comment\nbeta\ncafé\n"
	if err := os.WriteFile(input, []byte(content), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	dirs, err := ReadDirectories(input)
	if err != nil {
		t.Fatalf("ReadDirectories failed: %v", err)
	}

	if len(dirs) != 3 {
		t.Fatalf("expected 3 dirs, got %d", len(dirs))
	}
	if dirs[0] != "/alpha" || dirs[1] != "beta" || dirs[2] != "café" {
		t.Fatalf("unexpected dirs: %#v", dirs)
	}
}
