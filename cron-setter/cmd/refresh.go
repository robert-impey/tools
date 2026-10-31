package cmd

/*
Copyright © 2023 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

type refreshState int

const (
	stateBeforeHeader refreshState = iota
	stateInsideManagedBlock
	stateAfterFooter
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the cron-setter section of a crontab from stdin",
	Long: `Reads an existing crontab from stdin, preserves any lines before the
cron-setter header comment and after the footer comment, and regenerates
the cron-setter tasks in between.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return refreshCrontab(cmd.InOrStdin(), cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}

func refreshCrontab(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	state := stateBeforeHeader

	for scanner.Scan() {
		line := scanner.Text()

		switch state {
		case stateBeforeHeader:
			if isHeaderComment(line) {
				printAllTasks(w)
				state = stateInsideManagedBlock
			} else {
				fmt.Fprintln(w, line)
			}
		case stateInsideManagedBlock:
			if isFooterComment(line) {
				state = stateAfterFooter
			}
		case stateAfterFooter:
			fmt.Fprintln(w, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if state == stateBeforeHeader {
		printAllTasks(w)
	}

	return nil
}

func isHeaderComment(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), headerCommentPrefix)
}

func isFooterComment(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), footerCommentPrefix)
}
