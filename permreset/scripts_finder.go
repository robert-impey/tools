package main

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileHasShebang returns true if the file's first line starts with "#!"
func FindFilesWithShebang(root string) ([]string, error) {
	var results []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		has, err := FileHasShebang(path)
		if err != nil {
			return err
		}
		if has {
			results = append(results, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func FileHasShebang(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}
	line = strings.TrimRight(line, "\r\n")
	if line == "" {
		return false, nil
	}
	return strings.HasPrefix(line, "#!"), nil
}
