package lib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type FolderManager struct {
	Locations []string
	Folders   []string
}

// GetScriptPath joins the autoGenFolder with cleaned versions of location1 and location2.
// In Go, filepath.Join is the equivalent of Java's Path.resolve().
func GetScriptPath(autoGenFolder string, location1 string, location2 string) string {
	return filepath.Join(
		autoGenFolder,
		GetCleanLocationName(location1),
		GetCleanLocationName(location2),
	)
}

// GetCleanLocationName replaces special characters with underscores and trims trailing ones.
func GetCleanLocationName(location string) string {
	// Replaces one or more of ':', '\', '/', or ' ' with "_"
	reLegal := regexp.MustCompile(`[:\\/ ]+`)
	allLegal := reLegal.ReplaceAllString(location, "_")

	// Replaces one or more trailing underscores at the end of the string
	reTrailing := regexp.MustCompile(`_+$`)
	return reTrailing.ReplaceAllString(allLegal, "")
}

func (fm *FolderManager) ListManagedFolders(outFile io.Writer) {
	topOfPrintingLocations := true

	for _, location := range fm.Locations {
		// In Go, we usually just work with strings for paths
		if _, err := os.Stat(location); err != nil {
			// If location doesn't exist, skip (mimics Files.exists)
			continue
		}

		topOfPrintingFolders := true
		for _, folder := range fm.Folders {
			folderPath := filepath.Join(location, folder)

			// Lstat returns info about the link itself, not the target
			info, err := os.Lstat(folderPath)
			if err != nil {
				continue
			}

			// Check: Exists && IsDirectory && NOT a Symbolic Link
			if info.IsDir() && (info.Mode()&os.ModeSymlink == 0) {

				if topOfPrintingLocations {
					topOfPrintingLocations = false
				} else {
					if topOfPrintingFolders {
						fmt.Fprintln(outFile) // Print the blank line between locations
					}
				}

				absPath, err := filepath.Abs(folderPath)
				if err == nil {
					fmt.Fprintln(outFile, absPath)
				}

				topOfPrintingFolders = false
			}
		}
	}
}
