package internal_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robert-impey/tools/internal"
)

// writeFile is a small helper to create a temporary text file containing
// the provided lines.  It's used by several of the tests below.
func writeFile(t *testing.T, dir, name string, lines []string) string {
	t.Helper()
	scriptPath := filepath.Join(dir, name)
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(scriptPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file %s: %v", scriptPath, err)
	}
	return scriptPath
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestLoadFolderManager_SuccessfullyLoadsLocationsAndFolders(t *testing.T) {
	tempDir := t.TempDir()

	locationsFile := writeFile(t, tempDir, "locations.txt", []string{
		"C:/Data",
		"D:/Archive",
	})

	foldersFile := writeFile(t, tempDir, "folders.txt", []string{
		"logs",
		"configs",
	})

	manager, err := internal.LoadFolderManager(locationsFile, foldersFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if manager == nil {
		t.Fatalf("expected manager, got nil")
	}

	if !equal(manager.Locations, []string{"C:/Data", "D:/Archive"}) {
		t.Fatalf("unexpected locations: %#v", manager.Locations)
	}
	if !equal(manager.Folders, []string{"logs", "configs"}) {
		t.Fatalf("unexpected folders: %#v", manager.Folders)
	}
}

func TestLoadFolderManager_IgnoresBlankAndCommentLines(t *testing.T) {
	tempDir := t.TempDir()

	locationsFile := writeFile(t, tempDir, "locations.txt", []string{
		"",
		"   ",
		"# comment",
		"C:/RealPath",
	})

	foldersFile := writeFile(t, tempDir, "folders.txt", []string{
		"# header",
		"folderA",
		"",
		"folderB",
	})

	manager, err := internal.LoadFolderManager(locationsFile, foldersFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !equal(manager.Locations, []string{"C:/RealPath"}) {
		t.Fatalf("unexpected locations: %#v", manager.Locations)
	}
	if !equal(manager.Folders, []string{"folderA", "folderB"}) {
		t.Fatalf("unexpected folders: %#v", manager.Folders)
	}
}

func TestLoadFolderManager_ErrorsIfLocationsEmpty(t *testing.T) {
	tempDir := t.TempDir()

	locationsFile := writeFile(t, tempDir, "locations.txt", []string{
		"   ",
		"# comment",
	})

	foldersFile := writeFile(t, tempDir, "folders.txt", []string{"folderA"})

	_, err := internal.LoadFolderManager(locationsFile, foldersFile)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "locations cannot be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadFolderManager_ErrorsIfFoldersEmpty(t *testing.T) {
	tempDir := t.TempDir()

	locationsFile := writeFile(t, tempDir, "locations.txt", []string{"C:/Data"})

	foldersFile := writeFile(t, tempDir, "folders.txt", []string{
		"",
		"# comment",
	})

	_, err := internal.LoadFolderManager(locationsFile, foldersFile)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "folders cannot be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadFolderManager_ErrorsIfFileDoesNotExist(t *testing.T) {
	tempDir := t.TempDir()

	missing := filepath.Join(tempDir, "missing.txt")
	foldersFile := writeFile(t, tempDir, "folders.txt", []string{"folder"})

	_, err := internal.LoadFolderManager(missing, foldersFile)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to read locations") {
		t.Fatalf("expected error to mention missing file, got: %v", err)
	}
}
