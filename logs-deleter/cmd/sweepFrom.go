package cmd

/*
Copyright © 2022 Robert Impey, robert-impey@users.noreply.github.com
*/

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/robert-impey/tools/logs-deleter/lib"
	"github.com/spf13/cobra"
)

var Tool string

// sweepFromCmd represents the sweepFrom command
var sweepFromCmd = &cobra.Command{
	Use:   "sweepFrom",
	Short: "Sweep away the old log files for just one tool.",
	Long:  `Delete just the old log files for one tool that is logged.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := sweepFrom()
		if err != nil {
			log.Fatalln(err.Error())
		} else if Verbose {
			log.Println("Success")
		}
	},
}

func init() {
	rootCmd.AddCommand(sweepFromCmd)

	sweepFromCmd.Flags().StringVarP(&Tool, "tool", "t", "", "Tool to sweep")
}

func sweepFrom() error {
	if len(Tool) == 0 {
		return errors.New("tool not set")
	}

	var logsDir string
	if LogsDirectory == "" {
		log.Fatalln("LogsDirectory not set")
	} else {
		logsDir = LogsDirectory
	}

	var toolPath = filepath.Join(logsDir, Tool)

	_, err1 := os.Stat(toolPath)
	if err1 != nil {
		return err1
	}

	err2 := lib.DeleteFrom(toolPath, Days, DeleteEmpty, Verbose)
	if err2 != nil {
		return err2
	}
	return nil
}
