package internal

import (
	"bufio"
	"os"
	"strings"
)

// RsyncFileHasCopies reports whether fileName is an rsync log (i.e. it
// contains an incremental file list) and, if so, whether any files were copied.
func RsyncFileHasCopies(fileName string) (isRsyncLog bool, hasCopies bool, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return false, false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "receiving incremental file list" || line == "sending incremental file list" {
			isRsyncLog = true
			continue
		}
		if isRsyncLog && strings.HasPrefix(line, "sent ") {
			return true, false, nil
		}
		if isRsyncLog && line != "" {
			return true, true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, false, err
	}

	return isRsyncLog, false, nil
}
