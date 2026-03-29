package tdyfldr

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "tdyfldr",
	Short: "Search directories for matching file stems",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func logsDirPtr(logsDir string) *string {
	if logsDir == "" {
		return nil
	}
	return &logsDir
}

func addLogsDirFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "logs-dir", "", "Directory where logs should be written")
}
