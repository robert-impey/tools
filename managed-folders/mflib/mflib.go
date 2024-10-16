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

	commonLocalScriptsDirectory, err1 := GetCommonLocalScriptsDirectory()
	if err1 != nil {
		return "", err1
	}

	foldersFileName := "folders.txt"
	foldersFile := path.Join(commonLocalScriptsDirectory, foldersFileName)
	_, err2 := os.Stat(foldersFile)

	if err2 != nil {
		return "", err2
	}

	foldersFileAbs, err3 := filepath.Abs(foldersFile)

	if err3 != nil {
		return "", err3
	}
	return searchForFile(foldersFileName, foldersFileAbs)
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

	locationsFileName := "locations.txt"
	locationsFile := path.Join(commonLocalScriptsDirectory, osFolder, locationsFileName)
	_, err2 := os.Stat(locationsFile)

	if err2 != nil {
		return "", err2
	}

	locationsFileAbs, err3 := filepath.Abs(locationsFile)

	if err3 != nil {
		return "", err3
	}

	return searchForFile(locationsFileName, locationsFileAbs)
}

func searchForFile(fileName string, defaultFile string) (string, error) {
	localScriptsDirectory, err := GetLocalScriptsDirectory()
	if err != nil {
		return "", err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	machineLSDir := path.Join(localScriptsDirectory, hostname)

	_, err = os.Stat(machineLSDir)

	if os.IsNotExist(err) {
		return defaultFile, nil
	}

	userInfo, err := user.Current()
	if err != nil {
		return "", err
	}

	userMachineLSDir := path.Join(machineLSDir, userInfo.Username)

	userMachineFile := path.Join(userMachineLSDir, fileName)

	_, err = os.Stat(userMachineFile)
	if err == nil {
		absUserMachineFile, err := filepath.Abs(userMachineFile)
		if err != nil {
			log.Fatalln(err)
		}

		return absUserMachineFile, nil
	}

	machineFile := path.Join(machineLSDir, fileName)
	absMachineFile, err := filepath.Abs(machineFile)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = os.Stat(absMachineFile)

	if os.IsNotExist(err) {
		return defaultFile, nil
	}

	return absMachineFile, nil
}
