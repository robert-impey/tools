package internal

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	common "github.com/robert-impey/tools/internal"
)

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

	manager := &common.FolderManager{
		Locations: []string{loc1, loc2},
		Folders:   []string{"Shared"},
	}

	if err := GenerateRobocopyScripts(manager, autoGen); err != nil {
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

	manager := &common.FolderManager{
		Locations: []string{`C:\Data`},
		Folders:   []string{"Folder"},
	}

	err := GenerateRobocopyScripts(manager, missing)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "auto-generated folder does not exist") {
		t.Fatalf("expected error to contain %q, got %q",
			"Auto-generated folder does not exist", err.Error())
	}
}

func TestCreateRobocopySyncScript_GeneratesCorrectContent(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("robocopy script generation is Windows-only")
	}

	tempDir := t.TempDir()

	scriptPath := filepath.Join(tempDir, "script.ps1")
	sourcePath := `C:\Source`
	destinationPath := `D:\Destination`

	err := createRobocopySyncScript("MyFolder", scriptPath, sourcePath, destinationPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(scriptPath); err != nil {
		t.Fatalf("expected script file to exist: %v", err)
	}

	contentBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read script file: %v", err)
	}
	content := string(contentBytes)

	if !strings.Contains(content, `$folder = "MyFolder"`) {
		t.Errorf("expected folder assignment in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$src = "`+sourcePath+`"`) {
		t.Errorf("expected src assignment in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$dst = "`+destinationPath+`"`) {
		t.Errorf("expected dst assignment in script, got:\n%s", content)
	}

	if strings.Contains(content, `Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1"`) {
		t.Errorf("did not expect Synch module import in script, got:\n%s", content)
	}

	if strings.Contains(content, `function CleanFileName($fileName)`) {
		t.Errorf("did not expect CleanFileName helper in script, got:\n%s", content)
	}

	if strings.Contains(content, `function Get-LogsTimeStr {`) {
		t.Errorf("did not expect Get-LogsTimeStr helper in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$srcLogStr = "C_Source"`) {
		t.Errorf("expected literal source log name in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$dstLogStr = "D_Destination"`) {
		t.Errorf("expected literal destination log name in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$logTimeStr = Get-Date -Format "yyyy-MM-ddTHH_mm_ss"`) {
		t.Errorf("expected inline Get-LogsTimeStr expression in script, got:\n%s", content)
	}

	if !strings.Contains(content, `Start-Process ROBOCOPY`) {
		t.Errorf("expected robocopy invocation in script, got:\n%s", content)
	}
}

func TestCreateAllFoldersRobocopySyncScript_GeneratesCorrectContent(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("robocopy script generation is Windows-only")
	}

	tempDir := t.TempDir()

	scriptPath := filepath.Join(tempDir, "all.ps1")
	sourcePath := `C:\Source`
	destinationPath := `D:\Destination`
	folders := []string{"Folder1", "Folder2"}

	err := createAllFoldersRobocopySyncScript(folders, scriptPath, sourcePath, destinationPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(scriptPath); err != nil {
		t.Fatalf("expected script file to exist: %v", err)
	}

	contentBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read script file: %v", err)
	}
	content := string(contentBytes)

	if !strings.Contains(content, `$folders = "Folder1", "Folder2"`) {
		t.Errorf("expected folders assignment in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$src = "`+sourcePath+`"`) {
		t.Errorf("expected src assignment in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$dst = "`+destinationPath+`"`) {
		t.Errorf("expected dst assignment in script, got:\n%s", content)
	}

	if !strings.Contains(content, `foreach ($folder in $folders)`) {
		t.Errorf("expected foreach loop in script, got:\n%s", content)
	}

	if strings.Contains(content, `function CleanFileName($fileName)`) {
		t.Errorf("did not expect CleanFileName helper in script, got:\n%s", content)
	}

	if strings.Contains(content, `function Get-LogsTimeStr {`) {
		t.Errorf("did not expect Get-LogsTimeStr helper in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$srcLogStr = "C_Source"`) {
		t.Errorf("expected literal source log name in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$dstLogStr = "D_Destination"`) {
		t.Errorf("expected literal destination log name in script, got:\n%s", content)
	}

	if !strings.Contains(content, `$logTimeStr = Get-Date -Format "yyyy-MM-ddTHH_mm_ss"`) {
		t.Errorf("expected inline Get-LogsTimeStr expression in script, got:\n%s", content)
	}

	if strings.Contains(content, `Synch $folder $src $dst $logged`) {
		t.Errorf("did not expect Synch invocation in script, got:\n%s", content)
	}
}
