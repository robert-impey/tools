package mflib

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
)

func GetManagedFoldersFileName() (string, error) {
	theUser, err := user.Current()
	if err != nil {
		return "", err
	}

	managedFoldersFileName := path.Join(theUser.HomeDir, "autogen", "managed-folders.txt")

	return managedFoldersFileName, nil
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

func GetCommonLocalScriptsDirectory() (string, error) {
	localScriptsDirectory, err1 := GetLocalScriptsDirectory()
	if err1 != nil {
		return "", err1
	}

	commonDir := filepath.Join(localScriptsDirectory, "_Common")
	commonDirStat, err2 := os.Stat(commonDir)
	if err2 != nil {
		return "", err2
	}
	if !commonDirStat.IsDir() {
		return "", fmt.Errorf("%s is not a directory", commonDir)
	}
	absCommonLocalScriptsDir, err3 := filepath.Abs(commonDir)

	if err3 != nil {
		return "", err3
	}
	return absCommonLocalScriptsDir, nil
}

func GetFoldersFile() (string, error) {
	commonLocalScriptsDirectory, err1 := GetCommonLocalScriptsDirectory()
	if err1 != nil {
		return "", err1
	}

	foldersFile := path.Join(commonLocalScriptsDirectory, "folders.txt")
	_, err2 := os.Stat(foldersFile)

	if err2 != nil {
		return "", err2
	}

	foldersFileAbs, err3 := filepath.Abs(foldersFile)

	if err3 != nil {
		return "", err3
	}

	return foldersFileAbs, nil
}

func GetLocationsFile() (string, error) {
	commonLocalScriptsDirectory, err1 := GetCommonLocalScriptsDirectory()
	if err1 != nil {
		return "", err1
	}

	osFolder := "linux"
	if runtime.GOOS == "windows" {
		osFolder = "Windows"
	}

	locationsFile := path.Join(commonLocalScriptsDirectory, osFolder, "locations.txt")
	_, err2 := os.Stat(locationsFile)

	if err2 != nil {
		return "", err2
	}

	locationsFileAbs, err3 := filepath.Abs(locationsFile)

	if err3 != nil {
		return "", err3
	}

	return locationsFileAbs, nil
}
