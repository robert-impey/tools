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
	assert.Contains(t, s, "$sourceFolder = \"C:\\\"")
	assert.Contains(t, s, "$destinationFolder = \"D:\\\"")
	assert.Contains(t, s, "$logTimeStr = Get-Date -Format \"yyyy-MM-ddTHH_mm_ss\"")
	assert.Contains(t, s, "Start-Process ROBOCOPY -ArgumentList \"\"\"$($sourceFolder)\"\" \"\"$($destinationFolder)\"\" /xo \"\"$($file)\"\"\" `")
	assert.Contains(t, s, "file1.txt")
	assert.Contains(t, s, "file2.txt")
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
	assert.Contains(t, s, "$sourceFolder = \"C:\\\"")
	assert.Contains(t, s, "$destinationFolder = \"D:\\\"")
}
