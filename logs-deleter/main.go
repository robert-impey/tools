package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func main() {
	var err = sweepLogsDir()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	} else {
		fmt.Println("Success")
	}
}

func sweepLogsDir() error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	var logsDir = filepath.Join(user.HomeDir, "logs")

	if _, err := os.Stat(logsDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	subDirs, err := filepath.Glob(filepath.Join(logsDir, "*"))
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		subStat, err := os.Stat(subDir)
		if err != nil {
			return err
		}

		err = deleteFrom(filepath.Join(logsDir, subStat.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteFrom(subPath string) error {
	cutoff := time.Now().AddDate(0, -1, 0)

	fmt.Printf("Searching %v for files older than %v\n", subPath, cutoff)

	files, err := filepath.Glob(filepath.Join(subPath, "*"))
	if err != nil {
		return err
	}

	filesToDelete := make([]os.FileInfo, 0)
	for _, file := range files {
		fileStat, err := os.Stat(file)
		if err != nil {
			return err
		}
		if fileStat.ModTime().Before(cutoff) {
			filesToDelete = append(filesToDelete, fileStat)
		}
	}

	fmt.Printf("Found %d files to delete in %v\n", len(filesToDelete), subPath)

	for _, fileToDelete := range filesToDelete {
		err = os.RemoveAll(fileToDelete.Name())
		if err != nil {
			return err
		}
	}

	return nil
}
