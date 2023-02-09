/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
)

// sweepNightlyCmd represents the sweepNightly command
var sweepNightlyCmd = &cobra.Command{
	Use:   "sweepNightly",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sweepNightly()
	},
}

func init() {
	rootCmd.AddCommand(sweepNightlyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sweepNightlyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sweepNightlyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sweepNightly() {
	localScriptsDirectory, err := getLocalScriptsDirectory()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	fmt.Printf("localScriptsDirectory: %s\n", localScriptsDirectory)

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	machineLSDir := path.Join(localScriptsDirectory, hostname)

	machineLSDirStat, err := os.Stat(machineLSDir)
	if err != nil {
		panic(err)
	}

	if machineLSDirStat.IsDir() {
		absMachineLSDir, err := filepath.Abs(machineLSDir)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s exists\n", absMachineLSDir)
	}
	machineLSNightly := path.Join(machineLSDir, "staydeleted", "nightly.txt")

	_, err = os.Stat(machineLSNightly)
	if err != nil {
		panic(err)
	}

	absMachineLSNightly, err := filepath.Abs(machineLSNightly)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s exists\n", absMachineLSNightly)
}

func getLocalScriptsDirectory() (string, error) {
	if dir, ok := os.LookupEnv("LOCAL_SCRIPTS"); ok {
		dirStat, err := os.Stat(dir)
		if err != nil {
			return "", err
		}
		if !dirStat.IsDir() {
			return "", err
		}
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return "", err
		}

		return absDir, nil
	}
	return "", errors.New("LOCAL_SCRIPTS var not set!")
}
