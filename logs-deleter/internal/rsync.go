package internal

import (
	"bufio"
	"os"
	"strings"
)

func RsyncFileHasCopies(fileName string) (bool, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return false, err
	}
	defer file.Close()

	fileListStarted := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "receiving incremental file list" || line == "sending incremental file list" {
			fileListStarted = true
			continue
		}
		if fileListStarted && strings.HasPrefix(line, "sent ") {
			return false, nil
		}
		if fileListStarted && line != "" {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
