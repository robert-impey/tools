package main

import (
	"fmt"
	"log"
	"os"
)

// ApplyPerms sets file mode 0755 on each path. If a file is already 0755, it is left unchanged.
func ApplyPerms(files []string) error {
	var firstErr error
	for _, f := range files {
		log.Printf("File with shebang: %s", f)

		info, err := os.Stat(f)
		if err != nil {
			log.Printf("Failed to read permissions for %s: %v", f, err)
			if firstErr == nil {
				firstErr = fmt.Errorf("stat failed for %s: %w", f, err)
			}
			continue
		}

		currentPerm := info.Mode().Perm()
		if currentPerm == 0755 {
			log.Printf("Leaving permissions unchanged for %s: already %04o", f, currentPerm)
			continue
		}

		if err := os.Chmod(f, 0755); err != nil {
			log.Printf("Failed to set permissions on %s: %v", f, err)
			if firstErr == nil {
				firstErr = fmt.Errorf("chmod failed for %s: %w", f, err)
			}
			continue
		}

		log.Printf("Changed permissions on %s from %04o to 0755", f, currentPerm)
	}
	return firstErr
}
