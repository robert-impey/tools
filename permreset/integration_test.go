package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIntegration_ApplyPermsSets0755(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping integration test on windows")
	}

	dir := t.TempDir()
	p := filepath.Join(dir, "script.sh")
	content := "#!/usr/bin/env sh\necho hi\n"
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Ensure initial permissions are not executable for owner
	if err := os.Chmod(p, 0644); err != nil {
		t.Fatalf("chmod initial: %v", err)
	}

	files, err := FindFilesWithShebang(dir)
	if err != nil {
		t.Fatalf("find files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}

	if err := ApplyPerms(files); err != nil {
		t.Fatalf("apply perms: %v", err)
	}

	fi, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	perm := fi.Mode().Perm()
	if perm != 0755 {
		t.Fatalf("expected perms 0755, got %o", perm)
	}
}
