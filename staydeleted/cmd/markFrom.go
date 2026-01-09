package cmd

/*
Copyright © 2022 - 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/robert-impey/tools/staydeleted/sdlib"
	"github.com/spf13/cobra"
)

// markFromCmd represents the markFrom command
var markFromCmd = &cobra.Command{
	Use:   "markFrom",
	Short: "Mark all the files in a text file for deletion",
	Long:  `If many files need to be marked for deletion, a text file can be provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			err := markFrom(arg)
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("Success")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(markFromCmd)
}

func markFrom(markFromFileName string) error {
	log.Printf("Reading %v\n", markFromFileName)

	markFromFile, err := os.Open(markFromFileName)

	if err != nil {
		return err
	}
	defer markFromFile.Close()

	filesToMark := make([]string, 0)

	input := bufio.NewScanner(markFromFile)
	for input.Scan() {
		fileToMark := input.Text()
		if len(strings.TrimSpace(fileToMark)) == 0 {
			continue
		}
		if strings.HasPrefix(fileToMark, "#") {
			continue
		}

		filesToMark = append(filesToMark, fileToMark)
	}

	sdlib.MarkFiles(filesToMark, sdlib.Delete)

	return nil
}
