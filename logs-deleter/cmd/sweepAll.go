package cmd

/*
Copyright © 2022 Robert Impey, robert.impey@hotmail.co.uk
*/

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"robertimpey.com/tools/logs-deleter/lib"
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
	logsDir, err := lib.GetLogsDir()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	const toolName = "logs-deleter"
	toolLogDir := filepath.Join(logsDir, toolName)
	if _, err := os.Stat(toolLogDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(toolLogDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	}

	timeStr := time.Now().Format("2006-01-02_15.04.05")
	outLogFileName := filepath.Join(toolLogDir, fmt.Sprintf("%s.log", timeStr))
	errLogFileName := filepath.Join(toolLogDir, fmt.Sprintf("%s.err", timeStr))

	outLogFile, err := os.Create(outLogFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	errLogFile, err := os.Create(errLogFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	err = sweepLogsDir(logsDir, outLogFile)
	if err != nil {
		fmt.Fprint(os.Stderr, errLogFile)
	} else if Verbose {
		fmt.Fprintf(outLogFile, "Success")
	}
}

func sweepLogsDir(logsDir string, outWriter io.Writer) error {
	subDirs, err := filepath.Glob(filepath.Join(logsDir, "*"))
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		subStat, err := os.Stat(subDir)
		if err != nil {
			return err
		}

		err = lib.DeleteFrom(filepath.Join(logsDir, subStat.Name()), Days, DeleteEmpty, outWriter, Verbose)
		if err != nil {
			return err
		}
	}

	return nil
}
