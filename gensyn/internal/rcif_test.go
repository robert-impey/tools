package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSynchFile_ShouldReturnCorrectSourceAndDestinationPaths(t *testing.T) {
	// Path to the existing test data in the repository
	rel := filepath.Join("test_data", "robocopy", "files", "fruit.txt")
	path := filepath.Clean(rel)

	sf, err := ParseSynchFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, sf)

	assert.Equal(t, "fruit", sf.Id)
	assert.Equal(t, `C:\`, sf.Source)
	assert.Equal(t, `D:\`, sf.Destination)
	assert.Equal(t, 3, len(sf.Files))
	assert.Contains(t, sf.Files, "apples.txt")
	assert.Contains(t, sf.Files, "bananas.docx")
	assert.Contains(t, sf.Files, "cherries.pdf")
}

func TestParseSynchFile_ShouldErrorForBadFiles(t *testing.T) {
	cases := []string{
		filepath.Join("test_data", "robocopy", "files", "no_files.txt"),
		filepath.Join("test_data", "robocopy", "files", "no_blank_line.txt"),
		filepath.Join("test_data", "empty.txt"),
	}

	for _, c := range cases {
		_, err := ParseSynchFile(filepath.Clean(c))
		assert.NotNil(t, err)
	}
}

func TestGenerateRcifScript_ShouldCreateScriptFile_WithExpectedContent(t *testing.T) {
	tempDir := t.TempDir()
	// create a simple synch file
	sfPath := filepath.Join(tempDir, "mysynch.txt")
	content := strings.Join([]string{`C:\`, `D:\`, "", "file1.txt", "file2.txt"}, "\n")
	err := os.WriteFile(sfPath, []byte(content), 0o644)
	assert.Nil(t, err)

	autogen := t.TempDir()
	scriptName := "myscript"

	err = GenerateRcifScript(sfPath, autogen, scriptName)
	assert.Nil(t, err)

	out := filepath.Join(autogen, scriptName+".ps1")
	b, err := os.ReadFile(out)
	assert.Nil(t, err)
	s := string(b)
	assert.Contains(t, s, "# AUTOGEN'D - DO NOT EDIT!")
	assert.NotContains(t, s, "Import-Module \"$($env:LOCAL_SCRIPTS)\\_Common\\synch\\Synch.psm1\"")
	assert.NotContains(t, s, "SynchSingleFile2Ways")
	assert.NotContains(t, s, "$sourceFolder = \"C:\\\"")
	assert.NotContains(t, s, "$destinationFolder = \"D:\\\"")
	assert.NotContains(t, s, "$file = \"file1.txt\"")
	assert.Contains(t, s, "$logTimeStr = Get-Date -Format \"yyyy-MM-ddTHH_mm_ss\"")
	assert.Contains(t, s, `Start-Process ROBOCOPY -ArgumentList @("C:\", "D:\", "/xo", "file1.txt") `)
	assert.Contains(t, s, `Start-Process ROBOCOPY -ArgumentList @("C:\", "D:\", "/xo", "file2.txt") `)
	assert.Contains(t, s, `Start-Process ROBOCOPY -ArgumentList @("D:\", "C:\", "/xo", "file1.txt") `)
	assert.Contains(t, s, "file1.txt")
	assert.Contains(t, s, "file2.txt")
}

func TestGenerateRcifScript_ShouldSynchronizeBothDirectionsWhenLogged(t *testing.T) {
	tempDir := t.TempDir()
	sfPath := filepath.Join(tempDir, "mysynch.txt")
	content := strings.Join([]string{`C:\`, `D:\`, "", "file1.txt"}, "\n")
	err := os.WriteFile(sfPath, []byte(content), 0o644)
	assert.Nil(t, err)

	autogen := t.TempDir()
	err = GenerateRcifScript(sfPath, autogen, "myscript")
	assert.Nil(t, err)

	out := filepath.Join(autogen, "myscript.ps1")
	b, err := os.ReadFile(out)
	assert.Nil(t, err)
	s := string(b)

	sourceToDestination := `Start-Process ROBOCOPY -ArgumentList @("C:\", "D:\", "/xo", "file1.txt")`
	destinationToSource := `Start-Process ROBOCOPY -ArgumentList @("D:\", "C:\", "/xo", "file1.txt")`
	assert.Equal(t, 2, strings.Count(s, sourceToDestination))
	assert.Equal(t, 2, strings.Count(s, destinationToSource))
	assert.Equal(t, 2, strings.Count(s, "-RedirectStandardOutput"))
	assert.Equal(t, 2, strings.Count(s, "-RedirectStandardError"))
}

func TestGenerateRcifScript_ShouldOverwriteExistingScriptFile(t *testing.T) {
	tempDir := t.TempDir()
	sfPath := filepath.Join(tempDir, "mysynch.txt")
	content := strings.Join([]string{`C:\`, `D:\`, "", "file1.txt"}, "\n")
	err := os.WriteFile(sfPath, []byte(content), 0o644)
	assert.Nil(t, err)

	autogen := t.TempDir()
	scriptName := "myscript"
	out := filepath.Join(autogen, scriptName+".ps1")

	// create an existing file
	err = os.WriteFile(out, []byte("Old content"), 0o644)
	assert.Nil(t, err)

	err = GenerateRcifScript(sfPath, autogen, scriptName)
	assert.Nil(t, err)

	b, err := os.ReadFile(out)
	assert.Nil(t, err)
	s := string(b)
	assert.NotContains(t, s, "Old content")
	assert.NotContains(t, s, "SynchSingleFile2Ways")
	assert.NotContains(t, s, "$sourceFolder = \"C:\\\"")
	assert.NotContains(t, s, "$destinationFolder = \"D:\\\"")
	assert.Contains(t, s, `Start-Process ROBOCOPY -ArgumentList @("C:\", "D:\", "/xo", "file1.txt") `)
}

func TestGenerateRcifScript_ShouldValidateRequiredArguments(t *testing.T) {
	tempDir := t.TempDir()
	synchFile := filepath.Join(tempDir, "synch.txt")
	assert.NoError(t, os.WriteFile(synchFile, []byte(strings.Join([]string{`C:\`, `D:\`, "", "file.txt"}, "\n")), 0o644))

	tests := []struct {
		name          string
		synchFilePath string
		autoGenDir    string
		scriptName    string
		wantError     string
	}{
		{
			name:       "missing synch file path",
			autoGenDir: tempDir,
			scriptName: "script",
			wantError:  "synch file path is required",
		},
		{
			name:          "missing output directory",
			synchFilePath: synchFile,
			scriptName:    "script",
			wantError:     "autogen folder is required",
		},
		{
			name:          "missing script name",
			synchFilePath: synchFile,
			autoGenDir:    tempDir,
			wantError:     "script name is required",
		},
		{
			name:          "nonexistent output directory",
			synchFilePath: synchFile,
			autoGenDir:    filepath.Join(tempDir, "missing"),
			scriptName:    "script",
			wantError:     "auto-generated folder does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GenerateRcifScript(tt.synchFilePath, tt.autoGenDir, tt.scriptName)
			assert.ErrorContains(t, err, tt.wantError)
		})
	}
}

func TestGenerateRcifScript_ShouldReturnErrorForInvalidSynchFile(t *testing.T) {
	tempDir := t.TempDir()
	synchFile := filepath.Join(tempDir, "invalid.txt")
	assert.NoError(t, os.WriteFile(synchFile, []byte("invalid"), 0o644))

	err := GenerateRcifScript(synchFile, tempDir, "script")
	assert.Error(t, err)
	assert.ErrorContains(t, err, "at least four lines")
}

func TestGenerateRcifScript_ShouldSkipCommentLines(t *testing.T) {
	tempDir := t.TempDir()
	synchFile := filepath.Join(tempDir, "synch.txt")
	content := strings.Join([]string{`C:\`, `D:\`, "", "# this is a comment", "file.txt"}, "\n")
	assert.NoError(t, os.WriteFile(synchFile, []byte(content), 0o644))

	err := GenerateRcifScript(synchFile, tempDir, "script")
	assert.NoError(t, err)

	script, err := os.ReadFile(filepath.Join(tempDir, "script.ps1"))
	assert.NoError(t, err)
	assert.NotContains(t, string(script), "# this is a comment")
	assert.Contains(t, string(script), `Start-Process ROBOCOPY -ArgumentList @("C:\", "D:\", "/xo", "file.txt") `)
}

func TestGenerateRcifScript_ShouldReturnErrorWhenOutputPathCannotBeCreated(t *testing.T) {
	tempDir := t.TempDir()
	synchFile := filepath.Join(tempDir, "synch.txt")
	content := strings.Join([]string{`C:\`, `D:\`, "", "file.txt"}, "\n")
	assert.NoError(t, os.WriteFile(synchFile, []byte(content), 0o644))

	err := GenerateRcifScript(synchFile, tempDir, filepath.Join("missing", "script"))
	assert.ErrorContains(t, err, "failed to create script file")
}

func TestGenerateRcifScript_ShouldReturnErrorWhenExistingScriptCannotBeRemoved(t *testing.T) {
	tempDir := t.TempDir()
	synchFile := filepath.Join(tempDir, "synch.txt")
	content := strings.Join([]string{`C:\`, `D:\`, "", "file.txt"}, "\n")
	assert.NoError(t, os.WriteFile(synchFile, []byte(content), 0o644))

	outputPath := filepath.Join(tempDir, "script.ps1")
	assert.NoError(t, os.Mkdir(outputPath, 0o755))
	assert.NoError(t, os.WriteFile(filepath.Join(outputPath, "existing.txt"), []byte("keep"), 0o644))

	err := GenerateRcifScript(synchFile, tempDir, "script")
	assert.ErrorContains(t, err, "failed to remove existing script")
}
