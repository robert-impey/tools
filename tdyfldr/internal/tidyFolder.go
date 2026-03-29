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

func BuildDirsAndFiles(name string) (map[string][]DirEntry, error) {
	dirsAndFiles := make(map[string][]DirEntry)

	err := filepath.WalkDir(name, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		parentDir := filepath.Dir(path)

		dirsAndFiles[parentDir] = append(dirsAndFiles[parentDir], DirEntry{
			Path:    path,
			Name:    d.Name(),
			IsDir:   d.IsDir(),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirsAndFiles, nil
}

func FindMatchingStems(dirsAndFiles map[string][]DirEntry) [][2]DirEntry {
	var matchingStems [][2]DirEntry

	for _, files := range dirsAndFiles {
		for _, file := range files {
			fileStem, fileExt, ok := splitStemExt(file.Path)
			if !ok {
				continue
			}

			for _, otherFile := range files {
				otherStem, otherExt, ok := splitStemExt(otherFile.Path)
				if !ok {
					continue
				}

				if fileExt != otherExt {
					continue
				}
				if otherStem == fileStem {
					continue
				}
				if strings.HasPrefix(otherStem, fileStem) {
					matchingStems = append(matchingStems, [2]DirEntry{file, otherFile})
				}
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
	var logSink io.Writer = os.Stdout
	var errSink io.Writer = os.Stderr

	if logsDir != "" {
		if err := os.MkdirAll(logsDir, 0o755); err != nil {
			return err
		}

		safe := sanitizeForFileName(dir)
		timestamp := getLogTime()

		logPath := filepath.Join(logsDir, fmt.Sprintf("%s-search-%s.log", timestamp, safe))
		errPath := filepath.Join(logsDir, fmt.Sprintf("%s-search-%s.err", timestamp, safe))

		logFile, err := os.Create(logPath)
		if err != nil {
			return err
		}
		defer logFile.Close()
		logSink = logFile

		errFile, err := os.Create(errPath)
		if err != nil {
			return err
		}
		defer errFile.Close()
		errSink = errFile
	}

	funcErr := func() error {
		dirsAndFiles, err := BuildDirsAndFiles(dir)
		if err != nil {
			return err
		}

		matchingStems := FindMatchingStems(dirsAndFiles)
		return printMatchingStems(logSink, dir, matchingStems)
	}()

	if funcErr != nil {
		_, _ = fmt.Fprintf(errSink, "ERROR processing %s: %v\n", dir, funcErr)
		return nil
	}

	if logsDir != "" {
		fmt.Printf("OK: processed %s\n", dir)
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
