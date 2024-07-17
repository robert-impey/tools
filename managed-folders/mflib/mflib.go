package mflib

/*
Copyright © 2024 Robert Impey robert.impey@hotmail.co.uk
*/

import (
	"fmt"
	"log"
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
	localScriptsDirectory, err := GetLocalScriptsDirectory()
	if err != nil {
		return "", err
	}

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

	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	machineLSDir := path.Join(localScriptsDirectory, hostname)

	_, err = os.Stat(machineLSDir)

	if os.IsNotExist(err) {
		return foldersFileAbs, nil
	}

	userInfo, err := user.Current()
	if err != nil {
		return "", err
	}

	userMachineLSDir := path.Join(machineLSDir, userInfo.Username)

	userMachineFolders := path.Join(userMachineLSDir, "folders.txt")

	_, err = os.Stat(userMachineFolders)
	if err == nil {
		absUserMachineFolders, err := filepath.Abs(userMachineFolders)
		if err != nil {
			log.Fatalln(err)
		}

		return absUserMachineFolders, nil
	}

	machineLSFolders := path.Join(machineLSDir, "folders.txt")
	absMachineLSFolders, err := filepath.Abs(machineLSFolders)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = os.Stat(absMachineLSFolders)

	if os.IsNotExist(err) {
		return foldersFileAbs, nil
	}

	return absMachineLSFolders, nil
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
