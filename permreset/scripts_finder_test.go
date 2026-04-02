package main

import (
	"path/filepath"
	"testing"
)

func TestFileHasShebang_Positive(t *testing.T) {
	p := filepath.Join("testdata", "unix-script.sh")
	ok, err := FileHasShebang(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected %s to have a shebang", p)
	}
}

func TestFileHasShebang_Negative(t *testing.T) {
	cases := []string{"Empty.txt", "NotScript.txt", "WindowsOnly.ps1"}
	for _, c := range cases {
		p := filepath.Join("testdata", c)
		ok, err := FileHasShebang(p)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", c, err)
		}
		if ok {
			t.Fatalf("did not expect %s to have a shebang", c)
		}
	}
}
