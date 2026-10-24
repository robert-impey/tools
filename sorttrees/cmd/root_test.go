package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCommandWritesSortedOutput(t *testing.T) {
	fileName := filepath.Join(t.TempDir(), "tree.txt")
	if err := os.WriteFile(fileName, []byte("b\n    z\n    a\n\na\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	rootCmd := newRootCmdWithStreams(&stdout, &stderr)
	rootCmd.SetArgs([]string{fileName, "--SeparateTopLevel"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := "a\n\nb\n    a\n    z\n "
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestCommandInPlaceCreatesBackup(t *testing.T) {
	dir := t.TempDir()
	fileName := filepath.Join(dir, "tree.txt")
	original := []byte("b\n    z\n    a\n")
	if err := os.WriteFile(fileName, original, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	rootCmd := newRootCmdWithStreams(&stdout, &stderr)
	rootCmd.SetArgs([]string{fileName, "--InPlace"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "b\n    a\n    z\n" {
		t.Errorf("file contents = %q", got)
	}

	backup, err := os.ReadFile(fileName + ".bak")
	if err != nil {
		t.Fatal(err)
	}
	if string(backup) != string(original) {
		t.Errorf("backup contents = %q, want %q", backup, original)
	}
}

func TestCommandInPlaceSkipsAlreadySortedFile(t *testing.T) {
	fileName := filepath.Join(t.TempDir(), "tree.txt")
	original := []byte("a\n    b\n")
	if err := os.WriteFile(fileName, original, 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	rootCmd := newRootCmdWithStreams(&stdout, &stderr)
	rootCmd.SetArgs([]string{fileName, "--InPlace"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if stdout.String() != "In place sorting requested but already sorted.\n" {
		t.Errorf("stdout = %q", stdout.String())
	}
	if _, err := os.Stat(fileName + ".bak"); !os.IsNotExist(err) {
		t.Errorf("backup exists or failed unexpectedly: %v", err)
	}
}
