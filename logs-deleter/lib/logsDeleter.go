package lib

import (
	"errors"
	"fmt"
	"io"
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

func DeleteFrom(subPath string, days int, deleteEmpty bool, outWriter io.Writer, verbose bool) error {
	cutoff := time.Now().AddDate(0, 0, -1*days)

	if verbose {
		fmt.Fprintf(outWriter, "Searching %v for files older than %v\n", subPath, cutoff)
	}

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
		} else if deleteEmpty && fileStat.Size() == 0 {
			filesToDelete = append(filesToDelete, fileStat)
		}
	}

	countFilesToDelete := len(filesToDelete)

	if verbose || countFilesToDelete > 0 {
		fmt.Fprintf(outWriter, "Found %d files to delete in %v\n", len(filesToDelete), subPath)
	}

	for _, fileToDelete := range filesToDelete {
		var pathToDelete = filepath.Join(subPath, fileToDelete.Name())
		fmt.Fprintf(outWriter, "Deleting %v\n", pathToDelete)

		err = os.RemoveAll(pathToDelete)
		if err != nil {
			return err
		}
	}

	return nil
}
