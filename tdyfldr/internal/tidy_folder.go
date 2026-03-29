package tidy_folder

import (
	"io/fs"
	"path/filepath"
	"sync"
)

type DirEntry struct {
	Path    string
	Name    string
	IsDir   bool
	ModTime string // Using string for simplicity, can be time.Time
	Size    int64
}

// buildDirsAndFiles translates the Rust build_dirs_and_files function to Go.
// It walks the directory tree rooted at 'name' and returns a map where keys
// are directory paths and values are slices of DirEntry structs representing
// the files within that directory.
func BuildDirsAndFiles(name string) (map[string][]DirEntry, error) {
	dirsAndFiles := make(map[string][]DirEntry)
	var mutex sync.Mutex // To protect map access during concurrent walks

	err := filepath.Walk(name, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil // Skip directories
		}

		parentDir := filepath.Dir(path)

		fileEntry := DirEntry{
			Path:    path,
			Name:    info.Name(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"), // Format time as string
			Size:    info.Size(),
		}

		mutex.Lock()
		dirsAndFiles[parentDir] = append(dirsAndFiles[parentDir], fileEntry)
		mutex.Unlock()

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirsAndFiles, nil
}
