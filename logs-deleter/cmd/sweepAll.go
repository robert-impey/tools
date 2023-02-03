/*
Copyright © 2022 Robert Impey, robert.impey@hotmail.co.uk
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"robertimpey.com/tools/logs-deleter/lib"
	"time"
)

var Sleep int32

// sweepAllCmd represents the sweepAll command
var sweepAllCmd = &cobra.Command{
	Use:   "sweepAll",
	Short: "Sweep all the log directories",
	Long: `Sweep all the log directories.

Find files that are older than an expiry.
Optionally wait a random number of seconds before starting.
This command logs its output.`,
	Run: func(cmd *cobra.Command, args []string) {
		sweepLogsDirWithLogs()
	},
}

func init() {
	rootCmd.AddCommand(sweepAllCmd)
	sweepAllCmd.Flags().Int32VarP(&Sleep, "sleep", "s", 0,
		"The maximum number of seconds to sleep before starting. A random time during the period is chosen.")
}

func sweepLogsDirWithLogs() {
	if Sleep > 0 {
		wait := rand.Int31n(Sleep)
		time.Sleep(time.Duration(wait) * time.Second)
	}

	logsDir, err := lib.GetLogsDir()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		return
	}

	const toolName = "logs-deleter"
	timeStr := time.Now().Format("2006-01-02_15.04.05")
	outLogFileName := filepath.Join(logsDir, toolName, fmt.Sprintf("%s.log", timeStr))
	errLogFileName := filepath.Join(logsDir, toolName, fmt.Sprintf("%s.err", timeStr))

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
	} else {
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

		err = lib.DeleteFrom(filepath.Join(logsDir, subStat.Name()), Days, DeleteEmpty, outWriter)
		if err != nil {
			return err
		}
	}

	return nil
}
