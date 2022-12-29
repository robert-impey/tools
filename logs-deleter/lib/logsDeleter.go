package lib

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func GetLogsDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	var logsDir = filepath.Join(user.HomeDir, "logs")

	if _, err := os.Stat(logsDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return logsDir, nil
}

func DeleteFrom(subPath string, days int) error {
	cutoff := time.Now().AddDate(0, 0, -1*days)

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
		var pathToDelete = filepath.Join(subPath, fileToDelete.Name())
		fmt.Printf("Deleting %v\n", pathToDelete)

		err = os.RemoveAll(pathToDelete)
		if err != nil {
			return err
		}
	}

	return nil
}
