package lib

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// helper to create a file with a given size and mod time
func createFile(t *testing.T, dir, name string, size int, modTime time.Time) string {
	t.Helper()
	p := filepath.Join(dir, name)
	f, err := os.Create(p)
	if err != nil {
		t.Fatalf("create file: %v", err)
	}
	if size > 0 {
		if _, err := f.Write(make([]byte, size)); err != nil {
			_ = f.Close()
			t.Fatalf("write file: %v", err)
		}
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close file: %v", err)
	}
	// Set access and modification times
	if err := os.Chtimes(p, modTime, modTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	return p
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestDeleteFrom_DeletesOldFiles(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	days := 7
	now := time.Now()
	old := now.AddDate(0, 0, -(days + 1))
	recent := now.Add(-1 * time.Hour)

	oldFile := createFile(t, dir, "old.log", 10, old)
	recentFile := createFile(t, dir, "recent.log", 10, recent)

	if err := DeleteFrom(dir, days, false, false); err != nil {
		t.Fatalf("DeleteFrom error: %v", err)
	}

	if exists(oldFile) {
		t.Errorf("expected old file to be deleted: %s", oldFile)
	}
	if !exists(recentFile) {
		t.Errorf("expected recent file to remain: %s", recentFile)
	}
}

func TestDeleteFrom_DeletesEmptyWhenFlag(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	days := 30
	// Both files are recent (not old enough), but one is empty
	recent := time.Now().Add(-30 * time.Minute)

	emptyRecent := createFile(t, dir, "empty.txt", 0, recent)
	nonEmptyRecent := createFile(t, dir, "data.txt", 5, recent)

	if err := DeleteFrom(dir, days, true, false); err != nil {
		t.Fatalf("DeleteFrom error: %v", err)
	}

	if exists(emptyRecent) {
		t.Errorf("expected empty recent file to be deleted when deleteEmpty=true: %s", emptyRecent)
	}
	if !exists(nonEmptyRecent) {
		t.Errorf("expected non-empty recent file to remain: %s", nonEmptyRecent)
	}
}

func TestDeleteFrom_DoesNotDeleteRecentWhenFlagFalse(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	days := 10
	recent := time.Now().Add(-10 * time.Minute)

	emptyRecent := createFile(t, dir, "empty.txt", 0, recent)
	nonEmptyRecent := createFile(t, dir, "data.txt", 5, recent)

	if err := DeleteFrom(dir, days, false, false); err != nil {
		t.Fatalf("DeleteFrom error: %v", err)
	}

	if !exists(emptyRecent) {
		t.Errorf("expected empty recent file to remain when deleteEmpty=false: %s", emptyRecent)
	}
	if !exists(nonEmptyRecent) {
		t.Errorf("expected non-empty recent file to remain: %s", nonEmptyRecent)
	}
}
