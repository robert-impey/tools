package cmd

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/robert-impey/tools/gensyn/internal"
	"github.com/spf13/cobra"
)

// rsyncCmd represents the rsync command
var rsyncCmd = &cobra.Command{
	Use:   "rsync",
	Short: "A program for generating scripts for synchronising directories using rsync",
	Long: `A program for generating scripts for synchronising directories using rsync
`,
	Run: func(cmd *cobra.Command, args []string) {
		rsync(args)
	},
}

var files bool
var autoGenDir string

func init() {
	rootCmd.AddCommand(rsyncCmd)

	rsyncCmd.Flags().BoolVarP(&files, "files", "f", false, "Generate the scripts for files rather than directories")
	rsyncCmd.Flags().StringVarP(&autoGenDir, "autogenDir", "a", "", "autogenDir")
}

func rsync(args []string) {
	if len(args) != 1 {
		log.Fatalln(fmt.Errorf("expected 1 argument got %v", len(args)))
	}

	gssFile := args[0]

	if autoGenDir == "" {
		log.Fatalln("You must set the autogenDir parameter")
	}

	if _, err := os.Stat(autoGenDir); errors.Is(err, os.ErrNotExist) {
		log.Fatalln(err.Error())
	}

	err := internal.GenerateSynchScripts(files, autoGenDir, gssFile)
	if err != nil {
		log.Fatalf("Unable to generate the scripts for %v - %v\n", gssFile, err)
	}
}
