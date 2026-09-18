package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunRerlDeletesLogsWithoutCopies(t *testing.T) {
	dir := t.TempDir()

	// Create log files
	noCopiesPath := filepath.Join(dir, "no-copies.robocopy-synch.log")
	withCopiesPath := filepath.Join(dir, "with-copies.robocopy-synch.log")

	contentNoCopies := "header\n    Files :       5117          0       5117          0          0          0\nfooter\n"
	contentWithCopies := "header\n    Files :       5117          123       5117          0          0          0\nfooter\n"

	if err := os.WriteFile(noCopiesPath, []byte(contentNoCopies), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.WriteFile(withCopiesPath, []byte(contentWithCopies), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	oldLogsDirectory := LogsDirectory
	oldVerbose := Verbose
	LogsDirectory = dir
	Verbose = false
	t.Cleanup(func() {
		LogsDirectory = oldLogsDirectory
		Verbose = oldVerbose
	})

	if err := runRerl(); err != nil {
		t.Fatalf("runRerl error: %v", err)
	}

	if _, err := os.Stat(noCopiesPath); !os.IsNotExist(err) {
		t.Fatalf("expected log without copies to be deleted, stat error: %v", err)
	}
	if _, err := os.Stat(withCopiesPath); err != nil {
		t.Fatalf("expected log with copies to remain: %v", err)
	}
}

func TestRunRerlRequiresLogsDirectory(t *testing.T) {
	oldLogsDirectory := LogsDirectory
	LogsDirectory = "  "
	t.Cleanup(func() { LogsDirectory = oldLogsDirectory })

	if err := runRerl(); err == nil {
		t.Fatalf("expected runRerl to reject an empty logs directory")
	}
}
