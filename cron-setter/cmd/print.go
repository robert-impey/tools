package cmd

/*
Copyright © 2023 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"github.com/robert-impey/tools/cron-setter/internal"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print out a crontab",
	Long: `This command prints out a crontab with the daily tasks.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cronsetter.PrintTasks(cmd.OutOrStdout())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
