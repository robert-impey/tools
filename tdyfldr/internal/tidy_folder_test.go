package tidy_folder

import (
	"os"
	"path/filepath"
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

	err = os.Mkdir(dir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}
	err = os.Mkdir(dir2, 0755)
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
	if len(dirsAndFiles) != 3 { // tempDir, dir1, dir2
		t.Errorf("Expected 3 directories, got %d", len(dirsAndFiles))
	}

	// Check files in tempDir
	if files, ok := dirsAndFiles[tempDir]; ok {
		if len(files) != 1 {
			t.Errorf("Expected 1 file in %s, got %d", tempDir, len(files))
		}
		if !containsFile(files, "file4.txt") {
			t.Errorf("file4.txt not found in %s", tempDir)
		}
	} else {
		t.Errorf("%s not found in map", tempDir)
	}

	// Check files in dir1
	if files, ok := dirsAndFiles[dir1]; ok {
		if len(files) != 2 {
			t.Errorf("Expected 2 files in %s, got %d", dir1, len(files))
		}
		if !containsFile(files, "file1.txt") || !containsFile(files, "file2.txt") {
			t.Errorf("file1.txt or file2.txt not found in %s", dir1)
		}
	} else {
		t.Errorf("%s not found in map", dir1)
	}

	// Check files in dir2
	if files, ok := dirsAndFiles[dir2]; ok {
		if len(files) != 1 {
			t.Errorf("Expected 1 file in %s, got %d", dir2, len(files))
		}
		if !containsFile(files, "file3.txt") {
			t.Errorf("file3.txt not found in %s", dir2)
		}
	} else {
		t.Errorf("%s not found in map", dir2)
	}
}

// Helper function to create a file with content
func createFile(t *testing.T, path, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", path, err)
	}
}

// Helper function to check if a slice of DirEntry contains a file with a given name
func containsFile(files []DirEntry, name string) bool {
	for _, file := range files {
		if file.Name == name {
			return true
		}
	}
	return false
}
