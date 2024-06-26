/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/robert-impey/tools/managed-folders/mflib"
	"github.com/spf13/cobra"
	"log"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the environment",
	Long:  `The managed folders tool expects folders and files to be in conventional places.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := mflib.GetLocalScriptsDirectory()
		printVariable("Local Scripts Directory", dir, err)

		dir, err = mflib.GetCommonLocalScriptsDirectory()
		printVariable("Common Local Scripts Directory", dir, err)

		dir, err = mflib.GetFoldersFile()
		printVariable("Folders File", dir, err)

		dir, err = mflib.GetLocationsFile()
		printVariable("Locations File", dir, err)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func printVariable(name string, directory string, err error) {
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s: %s\n", name, directory)
}
