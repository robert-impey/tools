package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRsync(t *testing.T) {
	oldFiles, oldAutoGenDir := files, autoGenDir
	t.Cleanup(func() {
		files, autoGenDir = oldFiles, oldAutoGenDir
	})

	t.Run("generates scripts", func(t *testing.T) {
		autoGenDir = t.TempDir()
		files = false

		gssFile := filepath.Join("..", "internal", "test_data", "rsync", "qalhata.txt")
		if err := rsync([]string{gssFile}); err != nil {
			t.Fatalf("rsync returned an error: %v", err)
		}

		if _, err := os.Stat(filepath.Join(autoGenDir, "qalhata.ps1")); err != nil {
			t.Fatalf("expected generated script: %v", err)
		}
	})

	t.Run("returns error for missing output directory", func(t *testing.T) {
		autoGenDir = filepath.Join(t.TempDir(), "missing")

		err := rsync([]string{"config.gss"})
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected os.ErrNotExist, got %v", err)
		}
	})

	t.Run("returns error for missing config file", func(t *testing.T) {
		autoGenDir = t.TempDir()

		err := rsync([]string{"missing.gss"})
		if err == nil {
			t.Fatal("expected an error for a missing config file")
		}
	})
}
