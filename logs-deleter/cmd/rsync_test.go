package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunRsyncDeletesLogsWithoutCopies(t *testing.T) {
	dir := t.TempDir()
	emptyPath := copyRsyncFixture(t, dir, "empty.log")
	genScriptPath := copyRsyncFixture(t, dir, "gen-script.log")
	copiesPath := copyRsyncFixture(t, dir, "receiving-copies.log")

	oldLogsDirectory := LogsDirectory
	oldVerbose := Verbose
	LogsDirectory = dir
	Verbose = false
	t.Cleanup(func() {
		LogsDirectory = oldLogsDirectory
		Verbose = oldVerbose
	})

	if err := runRsync(); err != nil {
		t.Fatalf("runRsync error: %v", err)
	}

	if _, err := os.Stat(emptyPath); !os.IsNotExist(err) {
		t.Fatalf("expected empty rsync log to be deleted, stat error: %v", err)
	}

	if _, err := os.Stat(genScriptPath); err != nil {
		t.Fatalf("expected unrelated log file to remain: %v", err)
	}
	if _, err := os.Stat(copiesPath); err != nil {
		t.Fatalf("expected rsync log with copies to remain: %v", err)
	}
}

func TestRunRsyncRequiresLogsDirectory(t *testing.T) {
	oldLogsDirectory := LogsDirectory
	LogsDirectory = "  "
	t.Cleanup(func() { LogsDirectory = oldLogsDirectory })

	if err := runRsync(); err == nil {
		t.Fatalf("expected runRsync to reject an empty logs directory")
	}
}

func copyRsyncFixture(t *testing.T, dir, name string) string {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("..", "internal", "test_data", "rsync", name))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}
