/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		print()
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}

func print() {
	cron := `# m h  dom mon dow   command
0 * * * * /home/robert/executables/Linux/prod/x64/run-stay-deleted sweepNightly --sleep 3600
0 11 * * * /usr/bin/zsh /home/robert/local-scripts/_Common/reset-perms/reset-perms.sh
0 12 * * * /home/robert/executables/Linux/prod/x64/logs-deleter sweepAll --sleep 3600
0 13 * * * /usr/bin/zsh /home/robert/local-scripts/_Common/synch/run-nightly.sh`

	fmt.Println(cron)
}
