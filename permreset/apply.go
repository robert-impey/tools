package main

import (
	"fmt"
	"log"
	"os"
)

// ApplyPerms sets file mode 0755 on each path. Returns an error if any chmod fails.
func ApplyPerms(files []string) error {
	var firstErr error
	for _, f := range files {
		log.Printf("File with shebang: %s", f)
		if err := os.Chmod(f, 0755); err != nil {
			log.Printf("Failed to set permissions on %s: %v", f, err)
			if firstErr == nil {
				firstErr = fmt.Errorf("chmod failed for %s: %w", f, err)
			}
		}
	}
	return firstErr
}
