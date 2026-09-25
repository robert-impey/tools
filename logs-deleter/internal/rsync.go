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
	inFileList := false
	afterCommand := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "rsync ") {
			isRsyncLog = true
			afterCommand = true
		} else if line == "receiving incremental file list" || line == "sending incremental file list" {
			isRsyncLog = true
			inFileList = true
		} else if strings.HasPrefix(line, "sent ") {
			isRsyncLog = true
			inFileList = false
			afterCommand = false
		} else if (inFileList || afterCommand) && line != "" {
			hasCopies = true
		}
	}

	if err := scanner.Err(); err != nil {
		return false, false, err
	}

	return isRsyncLog, hasCopies, nil
}
