package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/robert-impey/tools/logs-deleter/internal"
	"github.com/spf13/cobra"
)

var rsyncCmd = &cobra.Command{
	Use:   "rsync",
	Short: "Remove empty rsync log files",
	Long: `Search the logs directory for rsync logs.
Files that do not contain any copied files are deleted.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runRsync()
		if err != nil {
			log.Fatalln(err.Error())
		} else if Verbose {
			log.Println("Success")
		}
	},
}

func init() {
	rootCmd.AddCommand(rsyncCmd)
}

func runRsync() error {
	if strings.TrimSpace(LogsDirectory) == "" {
		return errors.New("logsDirectory not set")
	}

	if _, err := os.Stat(LogsDirectory); err != nil {
		return err
	}

	matchingFiles, err := findRsyncLogFiles(LogsDirectory)
	if err != nil {
		return err
	}

	if Verbose {
		log.Printf("There are %d log files\n", len(matchingFiles))
	}

	for _, logFile := range matchingFiles {
		hasCopies, err := internal.RsyncFileHasCopies(logFile)
		if err != nil {
			log.Printf("Failed to process file %v: %v\n", logFile, err)
			continue
		}

		if hasCopies {
			log.Printf("%v has copies - keeping\n", logFile)
			continue
		}

		log.Printf("%v has no copies - deleting\n", logFile)
		if err := os.Remove(logFile); err != nil {
			return err
		}
	}

	return nil
}

func findRsyncLogFiles(root string) ([]string, error) {
	var matchingFiles []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".log") {
			matchingFiles = append(matchingFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return matchingFiles, nil
}
