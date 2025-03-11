package lib

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func DeleteFrom(subPath string, days int, deleteEmpty bool, verbose bool) error {
	cutoff := time.Now().AddDate(0, 0, -1*days)

	if verbose {
		log.Printf("Searching %v for files older than %v\n", subPath, cutoff)
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
		log.Printf("Found %d files to delete in %v\n", len(filesToDelete), subPath)
	}

	for _, fileToDelete := range filesToDelete {
		var pathToDelete = filepath.Join(subPath, fileToDelete.Name())
		log.Printf("Deleting %v\n", pathToDelete)

		err = os.RemoveAll(pathToDelete)
		if err != nil {
			return err
		}
	}

	return nil
}
