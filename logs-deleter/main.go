package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func main() {
	var err = sweepLogsDir()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	} else {
		fmt.Println("Success")
	}
}

func sweepLogsDir() error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	var logsDir = filepath.Join(user.HomeDir, "logs")

	if _, err := os.Stat(logsDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	subDirs, err := filepath.Glob(filepath.Join(logsDir, "*"))
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		subStat, err := os.Stat(subDir)
		if err != nil {
			return err
		}

		fmt.Printf("%v - %v", subStat.Name(), subStat.ModTime())
	}

	return nil
}
