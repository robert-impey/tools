package tdyfldr

import (
	"github.com/robert-impey/tools/tdyfldr/internal"
	"github.com/spf13/cobra"
)

var searchLogsDir string

var searchCmd = &cobra.Command{
	Use:   "search <directory>",
	Short: "Search a single directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var logsDirPtr *string
		if searchLogsDir != "" {
			logsDirPtr = &searchLogsDir
		}

		return internal.SearchDirectory(args[0], logsDirPtr)
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchLogsDir, "logs-dir", "", "Directory where logs should be written")
	rootCmd.AddCommand(searchCmd)
}
