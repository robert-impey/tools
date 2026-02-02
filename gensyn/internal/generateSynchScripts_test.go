package internal

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGSSFile(t *testing.T) {
	scriptsInfo, err := ParseGSSFile("merneith.txt")

	assert.Nil(t, err)
	assert.NotNil(t, scriptsInfo)

	assert.Equal(t, scriptsInfo.name, "merneith")
	assert.Equal(t, scriptsInfo.synch,
		"rsync --update --recursive --verbose --times --iconv=utf8 --perms --exclude-from=/home/robert/local-scripts/_Common/synch/rsync-excluded.txt")
	assert.Equal(t, scriptsInfo.src, "robert@merneith.robertimpey.com:~")
	assert.Equal(t, scriptsInfo.dst, "/home/robert")
	assert.Equal(t, len(scriptsInfo.items), 4)
}

func TestParseGSSFileBadFile(t *testing.T) {
	scriptsInfo, err := ParseGSSFile("does-not-exist.txt")

	assert.NotNil(t, err)
	assert.Nil(t, scriptsInfo)
}

func TestGenerateDirectorySynchScripts(t *testing.T) {
	outputDir := t.TempDir()

	err := GenerateSynchScripts(false, outputDir, "merneith.txt")
	assert.Nil(t, err)

	scriptsFile := path.Join(outputDir, "merneith.sh")
	if _, err := os.Stat(scriptsFile); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("scripts file %s doesn't exist", scriptsFile)
	}

	scriptsDir := path.Join(outputDir, "merneith")
	if _, err := os.Stat(scriptsDir); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("scripts dir %s doesn't exist", scriptsDir)
	}

	configScript := path.Join(scriptsDir, "config.sh")
	if _, err := os.Stat(configScript); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("config script file %s doesn't exist", configScript)
	}
}

func TestGenerateFilesSynchScripts(t *testing.T) {
	outputDir := t.TempDir()

	err := GenerateSynchScripts(true, outputDir, "ssh-config.txt")
	assert.Nil(t, err)

	scriptsFile := path.Join(outputDir, "ssh-config.sh")
	if _, err := os.Stat(scriptsFile); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("scripts file %s doesn't exist", scriptsFile)
	}
}

func TestGetCleanLocationName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes backslashes",
			input:    "C:\\Data",
			expected: "C_Data",
		},
		{
			name:     "removes forward slashes and trims trailing",
			input:    "D:/Archive/",
			expected: "D_Archive",
		},
		{
			name:     "removes spaces",
			input:    "E:\\My Documents",
			expected: "E_My_Documents",
		},
		{
			name:     "collapses multiple special characters",
			input:    "F::///  Data",
			expected: "F_Data",
		},
		{
			name:     "unix style paths",
			input:    "/var/data",
			expected: "_var_data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getCleanLocationName(tt.input)
			if actual != tt.expected {
				t.Errorf("getCleanLocationName(%q) = %q; want %q", tt.input, actual, tt.expected)
			}
		})
	}
}

func TestGetScriptPath(t *testing.T) {
	tests := []struct {
		name     string
		autoGen  string
		loc1     string
		loc2     string
		expected string
	}{
		{
			name:     "basic windows paths",
			autoGen:  `C:\AutoGen`,
			loc1:     `C:\Data`,
			loc2:     `D:\Backup`,
			expected: filepath.Join(`C:\AutoGen`, "C_Data", "D_Backup"),
		},
		{
			name:     "unix style paths",
			autoGen:  "/tmp/autogen",
			loc1:     "/var/data",
			loc2:     "/mnt/backup",
			expected: filepath.Join("/tmp/autogen", "_var_data", "_mnt_backup"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getScriptPath(tt.autoGen, tt.loc1, tt.loc2)
			if result != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestWriteScriptsCharacterization(t *testing.T) {
	info := &ScriptsInfo{
		name:  "merneith",
		synch: "rsync --flags",
		src:   "user@host:~",
		dst:   "/local/path",
		items: []string{"config", "data"},
	}

	t.Run("Directory Mode Full Output", func(t *testing.T) {
		outputDir := t.TempDir()
		err := writeScripts(false, outputDir, info)
		assert.Nil(t, err)

		// 1. Verify Main Script Content
		mainPath := filepath.Join(outputDir, "merneith.sh")
		content, _ := os.ReadFile(mainPath)
		script := string(content)

		// Verifying the 'To' and 'From' lines for the first item 'config'
		// This matches your dirsCmdLineTemplate logic
		assert.Contains(t, script, "rsync --flags user@host:~/config/ /local/path/config")
		assert.Contains(t, script, "rsync --flags /local/path/config/ user@host:~/config")

		// 2. Verify Sub-script creation
		// The code creates a folder named after the script for individual items
		subPath := filepath.Join(outputDir, "merneith", "config.sh")
		subContent, err := os.ReadFile(subPath)
		assert.Nil(t, err, "Sub-script should exist for 'config'")
		assert.Contains(t, string(subContent), "rsync --flags user@host:~/config/ /local/path/config")
	})

	t.Run("File Mode Full Output", func(t *testing.T) {
		outputDir := t.TempDir()
		fileInfo := &ScriptsInfo{
			name:  "dots",
			synch: "rsync --files-flags",
			src:   "/src/path",
			dst:   "/dst/path",
			items: []string{".bashrc"},
		}

		err := writeScripts(true, outputDir, fileInfo)
		assert.Nil(t, err)

		content, _ := os.ReadFile(filepath.Join(outputDir, "dots.sh"))
		script := string(content)

		// In file mode, item is only appended to the source
		assert.Contains(t, script, "rsync --files-flags /src/path/.bashrc /dst/path")
	})
}

func TestGenerateRobocopyScripts_CreatesExpectedScripts(t *testing.T) {
	tempDir := t.TempDir()

	loc1 := filepath.Join(tempDir, "Loc1")
	loc2 := filepath.Join(tempDir, "Loc2")

	// Create Shared subfolders
	if err := os.MkdirAll(filepath.Join(loc1, "Shared"), 0o755); err != nil {
		t.Fatalf("failed to create loc1/Shared: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(loc2, "Shared"), 0o755); err != nil {
		t.Fatalf("failed to create loc2/Shared: %v", err)
	}

	autoGen := filepath.Join(tempDir, "AutoGen")
	if err := os.MkdirAll(autoGen, 0o755); err != nil {
		t.Fatalf("failed to create AutoGen: %v", err)
	}

	manager := &FolderManager{
		Locations: []string{loc1, loc2},
		Folders:   []string{"Shared"},
	}

	if err := manager.GenerateRobocopyScripts(autoGen); err != nil {
		t.Fatalf("GenerateRobocopyScripts failed: %v", err)
	}

	// Forward direction
	scriptDir := getScriptPath(autoGen, loc1, loc2)

	if _, err := os.Stat(filepath.Join(scriptDir, "Shared.ps1")); err != nil {
		t.Errorf("expected Shared.ps1 to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(scriptDir, "_all.ps1")); err != nil {
		t.Errorf("expected _all.ps1 to exist: %v", err)
	}

	// Reverse direction
	reverseScriptDir := getScriptPath(autoGen, loc2, loc1)

	if _, err := os.Stat(filepath.Join(reverseScriptDir, "Shared.ps1")); err != nil {
		t.Errorf("expected reverse Shared.ps1 to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(reverseScriptDir, "_all.ps1")); err != nil {
		t.Errorf("expected reverse _all.ps1 to exist: %v", err)
	}
}

func TestGenerateRobocopyScripts_ErrorsIfAutoGenFolderDoesNotExist(t *testing.T) {
	tempDir := t.TempDir()

	missing := filepath.Join(tempDir, "MissingAutoGen")

	manager := &FolderManager{
		Locations: []string{`C:\Data`},
		Folders:   []string{"Folder"},
	}

	err := manager.GenerateRobocopyScripts(missing)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "auto-generated folder does not exist") {
		t.Fatalf("expected error to contain %q, got %q",
			"Auto-generated folder does not exist", err.Error())
	}
}

func writeFile(t *testing.T, dir, name string, lines []string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
	return path
}

func TestCreateFolderManager_SuccessfullyLoadsLocationsAndFolders(t *testing.T) {
	tempDir := t.TempDir()

	locationsFile := writeFile(t, tempDir, "locations.txt", []string{
		"C:/Data",
		"D:/Archive",
	})

	foldersFile := writeFile(t, tempDir, "folders.txt", []string{
		"logs",
		"configs",
	})

	manager, err := CreateFolderManager(locationsFile, foldersFile)
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

func TestCreateFolderManager_IgnoresBlankAndCommentLines(t *testing.T) {
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

	manager, err := CreateFolderManager(locationsFile, foldersFile)
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

func TestCreateFolderManager_ErrorsIfLocationsEmpty(t *testing.T) {
	tempDir := t.TempDir()

	locationsFile := writeFile(t, tempDir, "locations.txt", []string{
		"   ",
		"# comment",
	})

	foldersFile := writeFile(t, tempDir, "folders.txt", []string{"folderA"})

	_, err := CreateFolderManager(locationsFile, foldersFile)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "locations cannot be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateFolderManager_ErrorsIfFoldersEmpty(t *testing.T) {
	tempDir := t.TempDir()

	locationsFile := writeFile(t, tempDir, "locations.txt", []string{"C:/Data"})

	foldersFile := writeFile(t, tempDir, "folders.txt", []string{
		"",
		"# comment",
	})

	_, err := CreateFolderManager(locationsFile, foldersFile)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "folders cannot be empty" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateFolderManager_ErrorsIfFileDoesNotExist(t *testing.T) {
	tempDir := t.TempDir()

	missing := filepath.Join(tempDir, "missing.txt")
	foldersFile := writeFile(t, tempDir, "folders.txt", []string{"folder"})

	_, err := CreateFolderManager(missing, foldersFile)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to read locations") {
		t.Fatalf("expected error to mention missing file, got: %v", err)
	}
}

// helper
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
