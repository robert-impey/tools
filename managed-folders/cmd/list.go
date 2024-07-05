package cmd

/*
Copyright © 2024 Robert Impey robert.impey@hotmail.co.uk
*/

import (
	"bufio"
	"errors"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/robert-impey/tools/managed-folders/mflib"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the managed folders on this machine",
	Long:  `List the managed folders on this machine`,
	Run: func(cmd *cobra.Command, args []string) {
		folders, err1 := getFolders()
		locations, err2 := getLocations()

		err := errors.Join(err1, err2)
		if err != nil {
			log.Fatalln(err)
		}

		for i, location := range locations {
			havePrinted := false

			for _, folder := range folders {
				locatedFolder := path.Join(location, folder)
				locatedFolderStat, err := os.Stat(locatedFolder)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					log.Fatalln(err)
				}
				if !locatedFolderStat.IsDir() {
					log.Fatalf("%s is not a directory\n", locatedFolder)
				}
				absLocatedFolder, err := filepath.Abs(locatedFolder)
				if err != nil {
					log.Fatalln(err)
				}

				fmt.Println(absLocatedFolder)

				havePrinted = true
			}

			if havePrinted && i < len(locations)-1 {
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func getFolders() ([]string, error) {
	foldersFile, err := mflib.GetFoldersFile()
	if err != nil {
		return nil, err
	}

	return getDistinctSortedLine(foldersFile)
}

func getLocations() ([]string, error) {
	locationsFile, err := mflib.GetLocationsFile()
	if err != nil {
		return nil, err
	}

	return getDistinctSortedLine(locationsFile)
}

func getDistinctSortedLine(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	input := bufio.NewScanner(file)

	lines := mapset.NewSet[string]()
	for input.Scan() {
		line := strings.TrimRight(input.Text(), " /\\")
		if len(line) > 0 {
			lines.Add(line)
		}
	}

	linesSlice := lines.ToSlice()
	sort.Strings(linesSlice)

	return linesSlice, nil
}
