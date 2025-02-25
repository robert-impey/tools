package cmd

/*
Copyright © 2022-2025 Robert Impey, robert.impey@hotmail.co.uk
*/

import (
	"os"

	"github.com/spf13/cobra"
)

var LogsDirectory string
var Days int
var DeleteEmpty bool
var Verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "logs-deleter",
	Short: "A tool for deleting old log files",
	Long:  `A tool for deleting old log files`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&LogsDirectory, "logsDirectory", "l", "", "The logs directory")
	rootCmd.PersistentFlags().IntVarP(&Days, "days", "d", 30, "Days ago for cut off")
	rootCmd.PersistentFlags().BoolVarP(&DeleteEmpty, "deleteEmpty", "e", false, "Delete empty log files")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
}
