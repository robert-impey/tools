package cmd

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tools",
	Short: "A program for generating scripts for synchronising directories",
	Long: `A program for generating scripts for synchronising directories
`,
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
}
