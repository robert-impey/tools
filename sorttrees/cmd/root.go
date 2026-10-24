package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/robert-impey/tools/sorttrees/internal"
	"github.com/spf13/cobra"
)

func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	return newRootCmdWithStreams(os.Stdout, os.Stderr)
}

func newRootCmdWithStreams(stdout, stderr io.Writer) *cobra.Command {
	var separateTopLevel bool
	var inPlace bool

	rootCmd := &cobra.Command{
		Use:           "sorttrees TreeFile",
		Short:         "Sort an indented text tree",
		Long:          "Sort an indented text tree, treating each block at the same depth as a group.",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName := args[0]
			lines, err := internal.ReadLines(fileName)
			if err != nil {
				return err
			}

			if inPlace {
				alreadySorted, err := internal.AreLinesSortedTree(lines)
				if err != nil {
					return err
				}
				if alreadySorted {
					_, err := fmt.Fprintln(stdout, "In place sorting requested but already sorted.")
					return err
				}
				if err := internal.BackupFile(fileName); err != nil {
					return err
				}
			}

			output, err := internal.SortLines(lines, separateTopLevel)
			if err != nil {
				return err
			}
			if inPlace {
				return os.WriteFile(fileName, []byte(output), 0o666)
			}

			if _, err := io.WriteString(stdout, output); err != nil {
				return err
			}
			_, err = io.WriteString(stdout, " ")
			return err
		},
	}
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.Flags().BoolVar(&separateTopLevel, "SeparateTopLevel", false, "Separate top level trees with a blank line.")
	rootCmd.Flags().BoolVar(&inPlace, "InPlace", false, "Write the output on top of the input file.")

	return rootCmd
}
