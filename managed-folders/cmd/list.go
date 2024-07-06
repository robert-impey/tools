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
	"time"

	"github.com/spf13/cobra"
)

var write bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the managed folders on this machine",
	Long:  `List the managed folders on this machine`,
	Run: func(cmd *cobra.Command, args []string) {
		var output os.File

		if write {
			managedFoldersFile, err := mflib.GetManagedFoldersFile()
			if err != nil {
				log.Fatalln(err)
			}

			mfOut, err := os.Create(managedFoldersFile)

			output = *mfOut

			if err != nil {
				log.Fatalln(err)
			}
			fmt.Fprintln(&output, "# AUTOGEN'D - DO NOT EDIT!")
			
			now := time.Now().UTC()
			fmt.Fprintf(&output, "# Generated on %s\n\n", now.Format("2006-01-02"))
		} else {
			output = *os.Stdout
		}

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

				fmt.Fprintln(&output, absLocatedFolder)

				havePrinted = true
			}

			if havePrinted && i < len(locations)-1 {
				fmt.Fprintln(&output)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&write, "write", "w", false, "Write the managed folders file")
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
