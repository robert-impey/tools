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

// sweepAllCmd represents the sweepAll command
var sweepAllCmd = &cobra.Command{
	Use:   "sweepAll",
	Short: "Sweep all the log directories",
	Long: `Sweep all the log directories.

Find files that are older than an expiry.
This command logs its output.`,
	Run: func(cmd *cobra.Command, args []string) {
		sweepLogsDirWithLogs()
	},
}

func init() {
	rootCmd.AddCommand(sweepAllCmd)
}

func sweepLogsDirWithLogs() {
	var logsDir string
	if LogsDirectory == "" {
		log.Fatalln("LogsDirectory not set")
	} else {
		logsDir = LogsDirectory
	}

	const toolName = "logs-deleter"
	toolLogDir := filepath.Join(logsDir, toolName)
	if _, err := os.Stat(toolLogDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(toolLogDir, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	sweepErr := sweepLogsDir(logsDir)

	if sweepErr != nil {
		log.Fatalln(sweepErr)
	} else if Verbose {
		log.Println("Success")
	}
}

func sweepLogsDir(logsDir string) error {
	subDirs, err := filepath.Glob(filepath.Join(logsDir, "*"))
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		subStat, err := os.Stat(subDir)
		if err != nil {
			return err
		}

		err = lib.DeleteFrom(filepath.Join(logsDir, subStat.Name()), Days, DeleteEmpty, Verbose)
		if err != nil {
			return err
		}
	}

	return nil
}
