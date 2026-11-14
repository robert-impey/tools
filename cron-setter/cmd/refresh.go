package cmd

/*
Copyright © 2023 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"github.com/robert-impey/tools/cron-setter/internal"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the cron-setter section of a crontab from stdin",
	Long: `Reads an existing crontab from stdin, preserves any lines before the
cron-setter header comment and after the footer comment, and regenerates
the cron-setter tasks in between.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cronsetter.RefreshCrontab(cmd.InOrStdin(), cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
