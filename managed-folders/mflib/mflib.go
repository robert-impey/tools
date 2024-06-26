package mflib

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
)

func GetManagedFoldersFile() (string, error) {
	theUser, err := user.Current()
	if err != nil {
		return "", err
	}

	managedFoldersFile := path.Join(theUser.HomeDir, "autogen", "managed-folders.txt")

	_, err = os.Stat(managedFoldersFile)
	if err != nil {
		return "", err
	}

	return filepath.Abs(managedFoldersFile)
}

func GetLocalScriptsDirectory() (string, error) {
	if dir, ok := os.LookupEnv("LOCAL_SCRIPTS"); ok {
		dirStat, err := os.Stat(dir)
		if err != nil {
			return "", err
		}
		if !dirStat.IsDir() {
			return "", err
		}
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return "", err
		}

		return absDir, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	localScriptsDir := filepath.Join(currentUser.HomeDir, "local-scripts")
	localScriptsDirStat, err := os.Stat(localScriptsDir)
	if err != nil {
		return "", err
	}

	if !localScriptsDirStat.IsDir() {
		return "", err
	}
	absLocalScriptsDir, err := filepath.Abs(localScriptsDir)
	if err != nil {
		return "", err
	}

	return absLocalScriptsDir, nil
}
