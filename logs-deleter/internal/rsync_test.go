package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRsyncFileHasCopies_WhenCopiesExist(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "copies.log")
	content := "receiving incremental file list\npath/to/copied-file\n\nsent 164 bytes  received 21 bytes\n"
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	hasCopies, err := RsyncFileHasCopies(filePath)
	if err != nil {
		t.Fatalf("RsyncFileHasCopies error: %v", err)
	}
	if !hasCopies {
		t.Fatalf("expected RsyncFileHasCopies to return true")
	}
}

func TestRsyncFileHasCopies_WhenNoCopiesExist(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "empty.log")
	content := "receiving incremental file list\n\nsent 832 bytes  received 22,577 bytes\n"
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	hasCopies, err := RsyncFileHasCopies(filePath)
	if err != nil {
		t.Fatalf("RsyncFileHasCopies error: %v", err)
	}
	if hasCopies {
		t.Fatalf("expected RsyncFileHasCopies to return false")
	}
}

func TestRsyncFileHasCopies_WhenSummaryContainsNoFileList(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "summary.log")
	content := "receiving incremental file list\nsent 123 bytes  received 456 bytes\n"
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	hasCopies, err := RsyncFileHasCopies(filePath)
	if err != nil {
		t.Fatalf("RsyncFileHasCopies error: %v", err)
	}
	if hasCopies {
		t.Fatalf("expected summary line not to count as a copied file")
	}
}
