package mflib

/*
Copyright © 2024 Robert Impey robert.impey@hotmail.co.uk
*/

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

func GetLogsDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	var logsDir = filepath.Join(currentUser.HomeDir, "logs")

	if _, err := os.Stat(logsDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return logsDir, nil
}
