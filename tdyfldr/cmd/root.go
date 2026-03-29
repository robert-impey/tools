package tdyfldr

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "tdyfldr",
	Short: "Search directories for matching file stems",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
