package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSortLines(t *testing.T) {
	tests := []struct {
		name             string
		lines            []string
		separateTopLevel bool
		want             string
	}{
		{
			name: "sorts nested trees and removes gaps",
			lines: []string{
				"2",
				"    2",
				"    1",
				"1",
				"    2",
				"    1",
			},
			want: "1\n    1\n    2\n2\n    1\n    2\n",
		},
		{
			name: "separates top-level trees",
			lines: []string{
				"2",
				"1",
			},
			separateTopLevel: true,
			want:             "1\n\n2\n",
		},
		{
			name: "supports an empty tree",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortLines(tt.lines, tt.separateTopLevel)
			if err != nil {
				t.Fatalf("SortLines() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("SortLines() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSortLinesRejectsSkippedIndentationDepth(t *testing.T) {
	_, err := SortLines([]string{"root", "        grandchild"}, false)
	if err == nil {
		t.Fatal("SortLines() error = nil, want skipped-depth error")
	}
}

func TestAreLinesSortedTree(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{name: "empty", want: true},
		{name: "sorted children", lines: []string{"root", "    a", "    b"}, want: true},
		{name: "unsorted children", lines: []string{"root", "    b", "    a"}, want: false},
		{name: "unsorted after depth drop", lines: []string{"root", "    a", "        child", "    b", "next", "    z", "    c"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AreLinesSortedTree(tt.lines)
			if err != nil {
				t.Fatalf("AreLinesSortedTree() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("AreLinesSortedTree() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestReadLines(t *testing.T) {
	fileName := filepath.Join(t.TempDir(), "tree.txt")
	if err := os.WriteFile(fileName, []byte("one  \r\n\r\n    two \t\r\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := ReadLines(fileName)
	if err != nil {
		t.Fatalf("ReadLines() error = %v", err)
	}
	want := []string{"one", "    two"}
	if len(got) != len(want) {
		t.Fatalf("ReadLines() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("ReadLines()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestBackupFile(t *testing.T) {
	dir := t.TempDir()
	fileName := filepath.Join(dir, "tree.txt")
	if err := os.WriteFile(fileName, []byte("tree\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	modTime := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	if err := os.Chtimes(fileName, modTime, modTime); err != nil {
		t.Fatal(err)
	}

	if err := BackupFile(fileName); err != nil {
		t.Fatalf("BackupFile() error = %v", err)
	}

	backupName := fileName + ".bak"
	contents, err := os.ReadFile(backupName)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "tree\n" {
		t.Errorf("backup contents = %q, want %q", contents, "tree\n")
	}
	info, err := os.Stat(backupName)
	if err != nil {
		t.Fatal(err)
	}
	if !info.ModTime().Equal(modTime) {
		t.Errorf("backup modification time = %v, want %v", info.ModTime(), modTime)
	}
}
