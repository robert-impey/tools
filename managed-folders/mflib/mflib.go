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
