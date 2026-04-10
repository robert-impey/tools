package internal

import (
	"os"
	"regexp"
	"strconv"
)

var filesLineRegex = regexp.MustCompile(`\s*Files\s*:\s+\d+\s+(\d+)\s+`)

func IsFilesCopiedLine(line string) bool {
	match := filesLineRegex.FindStringSubmatch(line)
	if match == nil || len(match) < 2 {
		return false
	}

	copiedCount, err := strconv.Atoi(match[1])
	if err != nil {
		return false
	}

	return copiedCount > 0
}

func FileHasCopies(fileName string) (bool, error) {
	lines, err := os.ReadFile(fileName)
	if err != nil {
		return false, err
	}

	// iterate backwards
	start := len(lines) - 1
	end := len(lines)
	for i := start; i >= 0; i-- {
		if lines[i] == '\n' {
			if IsFilesCopiedLine(string(lines[i+1 : end])) {
				return true, nil // found a match near the end
			}
			end = i
		}
	}

	if IsFilesCopiedLine(string(lines[0:end])) {
		return true, nil
	}

	return false, nil // no match found
}
