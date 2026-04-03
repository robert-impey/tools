package cmd

import (
	"fmt"

	"github.com/robert-impey/tools/gensyn/internal"
	"github.com/spf13/cobra"
)

var (
	rcifFiles   string
	rcifScript  string
	rcifAutogen string
)

var rcifCmd = &cobra.Command{
	Use:   "rcif",
	Short: "Generate robocopy SynchSingleFile scripts from a synch file",
	Long:  `Generate robocopy SynchSingleFile2Ways scripts from a single synch file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if rcifFiles == "" {
			return fmt.Errorf("files path is required (use -f or --files)")
		}
		if rcifAutogen == "" {
			return fmt.Errorf("auto-generated folder is required (use -a or --autogen)")
		}
		if rcifScript == "" {
			return fmt.Errorf("script name is required (use -s or --script)")
		}

		return internal.GenerateRcifScript(rcifFiles, rcifAutogen, rcifScript)
	},
}

func init() {
	rootCmd.AddCommand(rcifCmd)

	rcifCmd.Flags().StringVarP(&rcifFiles, "files", "f", "", "The synch file to read")
	rcifCmd.Flags().StringVarP(&rcifAutogen, "autogen", "a", "", "The folder for the autogen'd scripts")
	rcifCmd.Flags().StringVarP(&rcifScript, "script", "s", "", "The script filename (without extension)")

	_ = rcifCmd.MarkFlagRequired("files")
	_ = rcifCmd.MarkFlagRequired("autogen")
	_ = rcifCmd.MarkFlagRequired("script")
}
