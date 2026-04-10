package internal

import (
	"os"
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

	dir := t.TempDir()
	filePath := filepath.Join(dir, "with-copies.robocopy-synch.log")

	content := "some header\n    Files :       5117          123       5117          0          0          0\nfooter\n"
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

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

	dir := t.TempDir()
	filePath := filepath.Join(dir, "no-copies.robocopy-synch.log")

	content := "some header\n    Files :       5117          0       5117          0          0          0\nfooter\n"
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	hasCopies, err := FileHasCopies(filePath)
	if err != nil {
		t.Fatalf("FileHasCopies error: %v", err)
	}

	if hasCopies {
		t.Fatalf("expected FileHasCopies to return false")
	}
}
