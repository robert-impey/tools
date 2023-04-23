package cmd

/*
Copyright © 2023 Robert Impey robert.impey@hotmail.co.uk
*/

import (
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read the current crontab and extract managed tasks",
	Long: `The contents of the crontab are read.

The tasks that are managed by this tool are within delimiting comments.`,
	Run: func(cmd *cobra.Command, args []string) {
		read()
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}

func read() {

}
