package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// FolderManager holds the configuration for generating sync scripts or
// listing managed folders. The two consuming packages (gensyn and lmf)
// each have their own helper methods, but they share the core data
// structure and the logic for loading locations/folders from text
// files.  By centralizing both here we avoid duplicating the file
// parsing code.

type FolderManager struct {
	Locations []string
	Folders   []string
}

// LoadFolderManager creates a FolderManager from the two paths.  It
// reads each file, ignores blank and comment lines, and ensures that
// neither resulting slice is empty.
func LoadFolderManager(locationsPath, foldersPath string) (*FolderManager, error) {
	locations, err := readLinesFromPath(locationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read locations: %w", err)
	}
	if len(locations) == 0 {
		return nil, fmt.Errorf("locations cannot be empty")
	}

	folders, err := readLinesFromPath(foldersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read folders: %w", err)
	}
	if len(folders) == 0 {
		return nil, fmt.Errorf("folders cannot be empty")
	}

	return &FolderManager{
		Locations: locations,
		Folders:   folders,
	}, nil
}

// readLinesFromPath loads a text file and returns non-blank, non-comment
// lines.  Comments are lines that start with '#'; all whitespace is
// trimmed from each line.
func readLinesFromPath(path string) ([]string, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", path)
		}
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}
