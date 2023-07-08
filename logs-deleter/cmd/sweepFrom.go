/*
Copyright © 2022 Robert Impey, robert.impey@hotmail.co.uk
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"robertimpey.com/tools/logs-deleter/lib"
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
			fmt.Fprint(os.Stderr, err.Error())
		} else if Verbose {
			fmt.Println("Success")
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

	logsDir, err := lib.GetLogsDir()
	if err != nil {
		return err
	}

	var toolPath = filepath.Join(logsDir, Tool)

	_, err1 := os.Stat(toolPath)
	if err1 != nil {
		return err1
	}

	err2 := lib.DeleteFrom(toolPath, Days, DeleteEmpty, os.Stdout, Verbose)
	if err2 != nil {
		return err2
	}
	return nil
}
