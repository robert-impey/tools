package lib

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetCleanLocationName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"C:\\Data", "C_Data"},
		{"D:/Archive/", "D_Archive"},
		{"E:\\My Documents", "E_My_Documents"},
	}

	for _, tt := range tests {
		result := GetCleanLocationName(tt.input)
		if result != tt.expected {
			t.Errorf("GetCleanLocationName(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetScriptPath(t *testing.T) {
	autoGen := "C:\\AutoGen"
	// filepath.Join handles platform-specific separators automatically
	expected := filepath.Join(autoGen, "C_Data", "D_Backup")

	result := GetScriptPath(autoGen, "C:\\Data", "D:\\Backup")

	if result != expected {
		t.Errorf("GetScriptPath() = %q; want %q", result, expected)
	}
}

func TestListManagedFolders(t *testing.T) {
	// 1. Create Temp Directory (Auto-cleaned by Go)
	tempDir := t.TempDir()

	// 2. Setup Paths
	location1 := filepath.Join(tempDir, "loc1")
	location2 := filepath.Join(tempDir, "loc2")

	// 3. Create Valid Directories
	os.MkdirAll(filepath.Join(location1, "folder1"), 0755)
	os.MkdirAll(filepath.Join(location1, "folder2"), 0755)
	os.MkdirAll(filepath.Join(location2, "folder1"), 0755)

	// 4. Create a File posing as a folder (to be ignored)
	os.MkdirAll(location2, 0755)
	fileAsFolder := filepath.Join(location2, "folder2")
	os.WriteFile(fileAsFolder, []byte("i am a file"), 0644)

	// 5. Initialize Manager
	locations := []string{location1, location2}
	folders := []string{"folder1", "folder2", "nonExistent"}

	fm := &FolderManager{
		Locations: locations,
		Folders:   folders,
	}

	// 6. Capture Output (Using a bytes.Buffer as the StringWriter)
	var buf bytes.Buffer

	// When
	fm.ListManagedFolders(&buf)

	// Then
	output := buf.String()
	expectedPath1, _ := filepath.Abs(filepath.Join(location1, "folder1"))
	expectedPath2, _ := filepath.Abs(filepath.Join(location1, "folder2"))
	expectedPath3, _ := filepath.Abs(filepath.Join(location2, "folder1"))
	invalidPath, _ := filepath.Abs(fileAsFolder)

	// Assertions
	if !strings.Contains(output, expectedPath1) {
		t.Errorf("Output should contain loc1/folder1, got: %s", output)
	}
	if !strings.Contains(output, expectedPath2) {
		t.Errorf("Output should contain loc1/folder2")
	}
	if !strings.Contains(output, expectedPath3) {
		t.Errorf("Output should contain loc2/folder1")
	}

	if strings.Contains(output, "nonExistent") {
		t.Error("Output should not contain non-existent folders")
	}
	if strings.Contains(output, invalidPath) {
		t.Error("Output should not contain files that match folder names")
	}
}

func TestListManagedFolders_ShouldSeparateLocationsWithEmptyLine(t *testing.T) {
	// 1. Setup Temp Directory
	tempDir := t.TempDir()

	location1 := filepath.Join(tempDir, "loc1")
	location2 := filepath.Join(tempDir, "loc2")

	// 2. Create directories to be found
	path1 := filepath.Join(location1, "folderA")
	path2 := filepath.Join(location2, "folderA")

	os.MkdirAll(path1, 0755)
	os.MkdirAll(path2, 0755)

	// 3. Initialize Manager
	// We use filepath.Abs to match the behavior of toAbsolutePath()
	absLoc1, _ := filepath.Abs(location1)
	absLoc2, _ := filepath.Abs(location2)

	fm := &FolderManager{
		Locations: []string{absLoc1, absLoc2},
		Folders:   []string{"folderA"},
	}

	// 4. Capture output
	var buf bytes.Buffer
	fm.ListManagedFolders(&buf)

	// 5. Assertions
	output := strings.TrimSpace(buf.String())
	// Handle cross-platform line endings by normalizing to \n
	normalizedOutput := strings.ReplaceAll(output, "\r\n", "\n")
	lines := strings.Split(normalizedOutput, "\n")

	// Expecting:
	// 0: Path to loc1/folderA
	// 1: (Empty String)
	// 2: Path to loc2/folderA
	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines (path, empty, path), but got %d:\n%q", len(lines), lines)
	}

	expectedPath1, _ := filepath.Abs(path1)
	expectedPath2, _ := filepath.Abs(path2)

	if lines[0] != expectedPath1 {
		t.Errorf("Line 0 mismatch.\nGot:  %s\nWant: %s", lines[0], expectedPath1)
	}

	if lines[1] != "" {
		t.Errorf("Line 1 should be empty as a separator, but got: %q", lines[1])
	}

	if lines[2] != expectedPath2 {
		t.Errorf("Line 2 mismatch.\nGot:  %s\nWant: %s", lines[2], expectedPath2)
	}
}

func TestListManagedFolders_ShouldIgnoreSymbolicLinks(t *testing.T) {
	// 1. Setup Temp Directory
	tempDir := t.TempDir()

	source := filepath.Join(tempDir, "real_source")
	location := filepath.Join(tempDir, "location")

	realFolder := filepath.Join(source, "myFolder")
	os.MkdirAll(realFolder, 0755)
	os.MkdirAll(location, 0755)

	// 2. Create a symlink in 'location' pointing to 'real_source/myFolder'
	symlinkPath := filepath.Join(location, "myFolder")
	err := os.Symlink(realFolder, symlinkPath)
	if err != nil {
		// Equivalent to catching UnsupportedOperationException or AccessDenied in Java
		t.Skipf("Skipping symlink test: operation not supported or permission denied: %v", err)
	}

	// 3. Initialize Manager
	fm := &FolderManager{
		Locations: []string{location},
		Folders:   []string{"myFolder"},
	}

	// 4. Capture output
	var buf bytes.Buffer

	// When
	fm.ListManagedFolders(&buf)

	// Then
	output := buf.String()
	if output != "" {
		t.Errorf("Output should be empty because symlinks must be ignored, but got: %q", output)
	}
}
