package internal

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/unicode/norm"
)

type DirEntry struct {
	Path    string
	Name    string
	IsDir   bool
	ModTime time.Time
	Size    int64
}

func BuildDirsAndFiles(name string) ([]DirEntry, error) {
	var entries []DirEntry

	err := filepath.WalkDir(name, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Log the error but continue traversal if possible, as this might be a transient I/O issue (like deadlock).
			fmt.Fprintf(os.Stderr, "Warning: Skipping path %s due to error during walkdir: %v\n", path, err)
			return nil // Return nil to signal WalkDir to continue with other paths
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			// Log the error but continue traversal if possible.
			fmt.Fprintf(os.Stderr, "Warning: Could not get info for %s: %v\n", path, err)
			return nil
		}

		entries = append(entries, DirEntry{
			Path:    path,
			Name:    d.Name(),
			IsDir:   d.IsDir(),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		})

		return nil
	})

	if err != nil {
		// If filepath.WalkDir returns an error here, it means the initial call failed or encountered a critical error not handled in the callback.
		return nil, fmt.Errorf("error walking directory %s: %w", name, err)
	}

	return entries, nil
}

func FindMatchingStems(dirsAndFiles []DirEntry) [][2]DirEntry {
	var matchingStems [][2]DirEntry

	for i := range dirsAndFiles {
		file := dirsAndFiles[i]
		fileStem, fileExt, ok := splitStemExt(file.Path)
		if !ok {
			continue
		}

		for j := range dirsAndFiles {
			otherFile := dirsAndFiles[j]
			otherStem, otherExt, ok := splitStemExt(otherFile.Path)
			if !ok {
				continue
			}

			// Ensure we are comparing two different files (i != j) and they have the same extension
			if i == j || fileExt != otherExt {
				continue
			}

			// Check for prefix relationship: Does 'otherStem' start with 'fileStem'?
			if strings.HasPrefix(otherStem, fileStem) && otherStem != fileStem {
				matchingStems = append(matchingStems, [2]DirEntry{file, otherFile})
			}
		}
	}

	return matchingStems
}

func ReadDirectories(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var dirs []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		dirs = append(dirs, norm.NFC.String(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dirs, nil
}

func SearchDirectory(dir string, logsDir string) error {
	dirsAndFiles, err := BuildDirsAndFiles(dir)
	if err != nil {
		return fmt.Errorf("failed to build directory structure for %s: %w", dir, err)
	}

	matchingStems := FindMatchingStems(dirsAndFiles)
	if len(matchingStems) == 0 {
		return nil
	}

	var out io.Writer = os.Stdout
	var logFile *os.File
	var errFile *os.File

	if logsDir != "" {
		safe := sanitizeForFileName(dir)
		timestamp := getLogTime()

		logPath := filepath.Join(logsDir, fmt.Sprintf("%s-search-%s.log", timestamp, safe))
		errPath := filepath.Join(logsDir, fmt.Sprintf("%s-search-%s.err", timestamp, safe))

		logFile, err = os.Create(logPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not create log file %s: %v\n", logPath, err)
		} else {
			defer logFile.Close()
			out = logFile
		}

		errFile, err = os.Create(errPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not create error log file %s: %v\n", errPath, err)
		} else {
			defer errFile.Close()
		}
	}

	funcErr := printMatchingStems(out, dir, matchingStems)

	if logsDir != "" {
		if logFile != nil && funcErr == nil {
			fmt.Printf("OK: processed %s\n", dir)
		}

		if funcErr != nil {
			if errFile != nil {
				_, _ = fmt.Fprintf(errFile, "ERROR processing %s: %v\n", dir, funcErr)
			}
			return fmt.Errorf("processing failed for %s: %w", dir, funcErr)
		}
	} else {
		if funcErr != nil {
			return fmt.Errorf("error processing %s: %w", dir, funcErr)
		}
	}

	return nil
}

func printMatchingStems(out io.Writer, name string, matchingStems [][2]DirEntry) error {
	if len(matchingStems) == 0 {
		return nil
	}

	if _, err := fmt.Fprintf(out, "\n--- Results for: %s ---\n", name); err != nil {
		return err
	}

	for _, pair := range matchingStems {
		file := pair[0]
		otherFile := pair[1]

		parent := filepath.Dir(file.Path)
		if _, err := fmt.Fprintln(out, parent); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "\t%s\n", file.Name); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "\t%s\n", otherFile.Name); err != nil {
			return err
		}
	}

	return nil
}

func getLogTime() string {
	return time.Now().Format("2006-01-02_15.04.05")
}

func sanitizeForFileName(s string) string {
	replacer := strings.NewReplacer("/", "_", `\`, "_", ":", "_")
	return replacer.Replace(s)
}

func splitStemExt(path string) (stem string, ext string, ok bool) {
	base := filepath.Base(path)
	if base == "." || base == string(filepath.Separator) {
		return "", "", false
	}

	ext = filepath.Ext(base)
	stem = strings.TrimSuffix(base, ext)
	return stem, ext, stem != ""
}
