package sdlib

import (
	"bufio"
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Action int

const (
	NoAction Action = iota
	Delete
	Keep
)

type ActionForFile struct {
	SdFile, File string
	Action       Action
}

type fileToDelete struct {
	Path, SDFile string
}

const SdFolderName = ".stay-deleted"

func GetActionForBool(keep bool) Action {
	if keep {
		return Keep
	}
	return Delete
}

func getStringForAction(action Action) string {
	if action == Keep {
		return "keep"
	}

	return "delete"
}

func getActionForString(actStr string) (Action, error) {
	if actStr == "delete" {
		return Delete, nil
	} else if actStr == "keep" {
		return Keep, nil
	}

	return NoAction, fmt.Errorf("unable to convert %s to an action", actStr)
}

func GetSdFolder(file string) (string, error) {
	dir := filepath.Dir(file)
	attemptedAbsSdFolder := filepath.Join(dir, SdFolderName)
	absSdFolder, err := filepath.Abs(attemptedAbsSdFolder)
	if err != nil {
		log.Printf("Unable to find the absolute path of '%v'!", attemptedAbsSdFolder)
		return "", err
	} else {
		return absSdFolder, nil
	}
}

func GetSdFile(file string) (string, error) {
	sdFolder, err := GetSdFolder(file)
	if err != nil {
		log.Printf("Unable to get sd folder for '%v'!", file)
		return "", err
	}

	fileBase := filepath.Base(file)
	data := []byte(fileBase)
	return filepath.Join(sdFolder, fmt.Sprintf("%x.txt", md5.Sum(data))), nil
}

func GetActionForFile(sdFileName, containingFolder string) (ActionForFile, error) {
	sdFile, err := os.Open(sdFileName)
	defer func(sdFile *os.File) {
		err := sdFile.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(sdFile)

	if err != nil {
		log.Printf("%v\n", err)
		return ActionForFile{"", "", NoAction}, err
	}

	input := bufio.NewScanner(sdFile)
	input.Scan()
	fileToProcessName := filepath.Join(containingFolder, input.Text())
	input.Scan()
	action, err := getActionForString(input.Text())
	if err != nil {
		log.Printf("%v\n", err)
		return ActionForFile{"", "", NoAction}, err
	}

	return ActionForFile{sdFileName, fileToProcessName, action}, nil
}

func SetActionForFile(fileName string, action Action) error {
	var absFileName, err = filepath.Abs(fileName)
	if err != nil {
		log.Printf("Unable to find the absolute path for '%v'!\n", fileName)
		return err
	}

	log.Printf("Marking: '%v'!\n", absFileName)
	fileBase := filepath.Base(absFileName)
	sdFileName, err := GetSdFile(absFileName)

	if err != nil {
		log.Printf("Unable to get sd file name for '%v'!",
			absFileName)
		return err
	}

	log.Printf("SD File: '%v'!\n", sdFileName)
	sdFolder := filepath.Dir(sdFileName)

	if _, err := os.Stat(sdFolder); os.IsNotExist(err) {
		log.Printf("Making directory '%v'\n", sdFolder)
		err := os.Mkdir(sdFolder, 0755)
		if err != nil {
			return err
		}
	}

	sdFile, err := os.Create(sdFileName)
	defer func(sdFile *os.File) {
		err := sdFile.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(sdFile)

	if err != nil {
		log.Printf("Couldn't create file '%v'!\n", sdFileName)
		return err
	}

	_, err = fmt.Fprintf(sdFile, "%v\n%s\n", fileBase, getStringForAction(action))
	if err != nil {
		return err
	}

	return nil
}

func ReadSweepFromFile(sweepFromFileName string) ([]string, error) {
	sweepFromFile, err := os.Open(sweepFromFileName)

	if err != nil {
		return nil, err
	}
	defer func(sdFile *os.File) {
		err := sdFile.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(sweepFromFile)
	directoriesToSweep := make([]string, 0)

	input := bufio.NewScanner(sweepFromFile)
	for input.Scan() {
		directoryToSweep := input.Text()
		if len(strings.TrimSpace(directoryToSweep)) == 0 {
			continue
		}
		if strings.HasPrefix(directoryToSweep, "#") {
			continue
		}

		directoriesToSweep = append(directoriesToSweep, directoryToSweep)
	}

	return directoriesToSweep, nil
}

func SweepFrom(sweepFromFileName string, expiryMonths int, verbose bool) error {
	var directoriesToSweepFrom, err = ReadSweepFromFile(sweepFromFileName)
	if err != nil {
		log.Printf("Unable to read file to sweep from '%v' - '%v'\n", sweepFromFileName, err)
		return err
	}

	for _, directoryToSweepFrom := range directoriesToSweepFrom {
		err := SweepDirectory(directoryToSweepFrom, expiryMonths, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

func SweepDirectory(directoryToSweep string, expiryMonths int, verbose bool) error {
	stat, err := os.Stat(directoryToSweep)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", directoryToSweep)
	}

	absDirectoryToSweep, err := filepath.Abs(directoryToSweep)
	if err != nil {
		log.Printf("Unable to find the absolute path for '%v' - '%v'!\n",
			directoryToSweep, err)
		return err
	}

	if verbose {
		fmt.Printf("Sweeping: '%v'\n", absDirectoryToSweep)
	}

	sdExpiryCutoff := time.Now().AddDate(0, -1*expiryMonths, 0)

	filesToDelete, err := findFilesToDelete(absDirectoryToSweep, sdExpiryCutoff, verbose)
	if err != nil {
		return err
	}

	deleteFilesToDelete(filesToDelete)

	return nil
}

func findFilesToDelete(absDirectoryToSweep string, sdExpiryCutoff time.Time, verbose bool) ([]fileToDelete, error) {
	filesToDelete := make([]fileToDelete, 0)
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}

		if info.IsDir() && info.Name() == SdFolderName {
			sdFolder := path
			if verbose {
				fmt.Printf("Search SD folder '%v'\n", sdFolder)
			}
			containingFolder := filepath.Dir(sdFolder)
			if verbose {
				fmt.Printf("Containing folder '%v'\n", containingFolder)
			}

			sdFiles, err := FindSdFiles(sdFolder)
			if err != nil {
				log.Printf("%v\n", err)
				return err
			}

			// Remove emptied sd folders
			if len(sdFiles) == 0 {
				log.Printf("Adding empty SD folder '%s' to the delete list\n", sdFolder)
				filesToDelete = append(filesToDelete, fileToDelete{Path: sdFolder, SDFile: ""})
			}

			for _, sdFile := range sdFiles {
				sdStat, err := os.Stat(sdFile)
				if err != nil {
					log.Printf("%v\n", err)
					return err
				}

				if !isSdFile(sdStat) {
					log.Printf("'%v' is not a legal name for SD file - deleting.\n",
						sdFile)
					filesToDelete = append(filesToDelete, fileToDelete{sdFile, ""})
					continue
				}

				if sdStat.ModTime().Before(sdExpiryCutoff) {
					log.Printf("Adding old SD file '%v' from %s to the delete list\n",
						sdFile,
						sdStat.ModTime().Format("2006-01-02 15:04:05"))
					filesToDelete = append(filesToDelete, fileToDelete{sdFile, ""})
					continue
				}

				if verbose {
					log.Printf("SD File '%v'\n", sdFile)
				}
				actionForFile, err := GetActionForFile(sdFile, containingFolder)
				if err != nil {
					log.Printf("%v\n", err)
					return err
				}

				if actionForFile.Action == Delete {
					if _, err := os.Stat(actionForFile.File); os.IsNotExist(err) {
						if verbose {
							log.Printf("'%v' already deleted.\n", actionForFile.File)
						}
						continue
					}
					log.Printf("Adding '%v' to the delete list\n", actionForFile.File)
					filesToDelete = append(filesToDelete, fileToDelete{actionForFile.File, actionForFile.SdFile})
				} else if actionForFile.Action == Keep {
					if verbose {
						log.Printf("Keeping '%v'\n", actionForFile.File)
					}
				} else {
					log.Printf("Unrecognised action '%v' from '%v'!\n",
						actionForFile.Action, sdFile)
					log.Printf("Adding unreadable SD file '%v' from %s to the delete list\n",
						sdFile,
						sdStat.ModTime().Format("2006-01-02 15:04:05"))
					filesToDelete = append(filesToDelete, fileToDelete{sdFile, ""})
				}
			}
		}

		return nil
	}

	err := filepath.Walk(absDirectoryToSweep, walker)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}

	return filesToDelete, nil
}

func deleteFilesToDelete(filesToDelete []fileToDelete) {
	var pe *fs.PathError
	for _, fileToDelete := range filesToDelete {
		var deleteMessage = fmt.Sprintf("Deleting '%v'", fileToDelete.Path)

		if len(fileToDelete.SDFile) > 0 {
			deleteMessage += fmt.Sprintf(" as instructed by '%v'", fileToDelete.SDFile)
		}
		log.Printf("%v\n", deleteMessage)

		err := os.RemoveAll(fileToDelete.Path)
		if err != nil {
			log.Printf("%v\n", err)
			if errors.As(err, &pe) {
				log.Printf("Failed to remove %v from %v\n", pe.Path, fileToDelete.SDFile)
			}
		}
	}
}

func FindSdFiles(sdFolder string) ([]string, error) {
	var sdFiles []string
	walker := func(sdFile string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		sdStat, err := os.Stat(sdFile)
		if err != nil {
			return err
		}
		if isSdFile(sdStat) {
			sdFiles = append(sdFiles, sdFile)
		}
		return nil
	}

	err := filepath.Walk(sdFolder, walker)

	if err != nil {
		return []string{}, err
	}
	return sdFiles, nil
}

func isSdFile(sdStat os.FileInfo) bool {
	re, _ := regexp.Compile(`[0-9a-fA-F]+.txt`)
	return re.Match([]byte(sdStat.Name()))
}
