package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunRerlDeletesLogsWithoutCopies(t *testing.T) {
	dir := t.TempDir()
	noCopiesPath := copyRerlFixture(t, dir, "empty.log")
	withCopiesPath := copyRerlFixture(t, dir, "copies.log")

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

func copyRerlFixture(t *testing.T, dir, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("..", "internal", "test_data", "robocopy", name))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	path := filepath.Join(dir, name+".robocopy-synch.log")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func TestRunRerlRequiresLogsDirectory(t *testing.T) {
	oldLogsDirectory := LogsDirectory
	LogsDirectory = "  "
	t.Cleanup(func() { LogsDirectory = oldLogsDirectory })

	if err := runRerl(); err == nil {
		t.Fatalf("expected runRerl to reject an empty logs directory")
	}
}
