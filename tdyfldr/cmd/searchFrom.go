package tdyfldr

import (
	"github.com/robert-impey/tools/tdyfldr/internal"
	"github.com/spf13/cobra"
)

var searchFromLogsDir string

var searchFromCmd = &cobra.Command{
	Use:   "search-from <directories-file>",
	Short: "Search multiple directories listed in a text file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var logsDirPtr *string
		if searchFromLogsDir != "" {
			logsDirPtr = &searchFromLogsDir
		}

		dirs, err := internal.ReadDirectories(args[0])
		if err != nil {
			return err
		}

		for _, dir := range dirs {
			if err := internal.SearchDirectory(dir, logsDirPtr); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	searchFromCmd.Flags().StringVar(&searchFromLogsDir, "logs-dir", "", "Directory where logs should be written")
	rootCmd.AddCommand(searchFromCmd)
}
