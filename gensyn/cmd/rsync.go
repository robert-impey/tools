/*
Copyright © 2025 Robert Impey robert.impey@hotmail.co.uk
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/robert-impey/tools/gensyn/lib"
	"github.com/spf13/cobra"
)

// rsyncCmd represents the rsync command
var rsyncCmd = &cobra.Command{
	Use:   "rsync",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rsync(args)
	},
}

var autoGenDir string

func init() {
	rootCmd.AddCommand(rsyncCmd)

	rsyncCmd.Flags().StringVarP(&autoGenDir, "autogenDir", "a", "", "autogenDir")
}

func rsync(args []string) {
	gssFiles := make([]string, 0)

	for _, arg := range args {
		_, err := os.Stat(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to process %v - %v\n", arg, err)
			continue
		}
		gssFiles = append(gssFiles, arg)
	}

	if autoGenDir == "" {
		log.Fatalln("You must set the autogenDir parameter")
	}

	if _, err := os.Stat(autoGenDir); errors.Is(err, os.ErrNotExist) {
		log.Fatalln(err.Error())
	}

	for _, gssFile := range gssFiles {
		err := lib.GenerateSynchScripts(autoGenDir, gssFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to generate the scripts for %v - %v\n", gssFile, err)
			continue
		}
	}
}
