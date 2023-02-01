/*
Copyright © 2022 Robert Impey, robert.impey@hotmail.co.uk
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"robertimpey.com/tools/logs-deleter/lib"
)

// sweepAllCmd represents the sweepAll command
var sweepAllCmd = &cobra.Command{
	Use:   "sweepAll",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sweepLogsDirWithLogs()
	},
}

func init() {
	rootCmd.AddCommand(sweepAllCmd)
}

func sweepLogsDirWithLogs() {
	var err = sweepLogsDir()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	} else {
		fmt.Println("Success")
	}
}

func sweepLogsDir() error {
	logsDir, err := lib.GetLogsDir()
	if err != nil {
		return err
	}

	subDirs, err := filepath.Glob(filepath.Join(logsDir, "*"))
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		subStat, err := os.Stat(subDir)
		if err != nil {
			return err
		}

		err = lib.DeleteFrom(filepath.Join(logsDir, subStat.Name()), Days, DeleteEmpty, os.Stdout)
		if err != nil {
			return err
		}
	}

	return nil
}
