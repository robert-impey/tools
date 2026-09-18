package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunRerlDeletesLogsWithoutCopies(t *testing.T) {
	dir := t.TempDir()
	// Load log files from fixed test data paths instead of generating temporary ones.
	const (
		testDataDir = "internal/test_data/robocopy"
	)

	// Paths to the scenario files:
	noCopiesPath := filepath.Join(testDataDir, "empty.log")
	withCopiesPath := filepath.Join(testDataDir, "copies.log")

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
