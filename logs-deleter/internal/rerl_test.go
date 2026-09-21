package internal

import (
	"path/filepath"
	"testing"
)

func TestIsFilesCopiedLine_NoFilesCopied(t *testing.T) {
	t.Parallel()

	line := "    Files :       5117          0       5117          0          0          0"

	if IsFilesCopiedLine(line) {
		t.Fatalf("expected IsFilesCopiedLine to return false")
	}
}

func TestIsFilesCopiedLine_FilesCopied(t *testing.T) {
	t.Parallel()

	line := "    Files :       5117          123       5117          0          0          0"

	if !IsFilesCopiedLine(line) {
		t.Fatalf("expected IsFilesCopiedLine to return true")
	}
}

func TestFileHasCopies_WhenCopiesExist(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join("test_data", "robocopy", "copies.log")

	hasCopies, err := FileHasCopies(filePath)

	if err != nil {
		t.Fatalf("FileHasCopies error: %v", err)
	}

	if !hasCopies {
		t.Fatalf("expected FileHasCopies to return true")
	}
}

func TestFileHasCopies_WhenNoCopiesExist(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join("test_data", "robocopy", "empty.log")
	hasCopies, err := FileHasCopies(filePath)

	if err != nil {
		t.Fatalf("FileHasCopies error: %v", err)
	}

	if hasCopies {
		t.Fatalf("expected FileHasCopies to return false")
	}
}
