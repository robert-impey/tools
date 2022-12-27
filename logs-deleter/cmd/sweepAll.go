/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// sweepAllCmd represents the sweepAll command
var sweepAllCmd = &cobra.Command{
	Use:   "sweepAll",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err = sweepLogsDir()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		} else {
			fmt.Println("Success")
		}
	},
}

func init() {
	rootCmd.AddCommand(sweepAllCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sweepAllCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sweepAllCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sweepLogsDir() error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	var logsDir = filepath.Join(user.HomeDir, "logs")

	if _, err := os.Stat(logsDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	subDirs, err := filepath.Glob(filepath.Join(logsDir, "*"))
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		subStat, err := os.Stat(subDir)
		if err != nil {
			return err
		}

		err = deleteFrom(filepath.Join(logsDir, subStat.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteFrom(subPath string) error {
	cutoff := time.Now().AddDate(0, -1, 0)

	fmt.Printf("Searching %v for files older than %v\n", subPath, cutoff)

	files, err := filepath.Glob(filepath.Join(subPath, "*"))
	if err != nil {
		return err
	}

	filesToDelete := make([]os.FileInfo, 0)
	for _, file := range files {
		fileStat, err := os.Stat(file)
		if err != nil {
			return err
		}
		if fileStat.ModTime().Before(cutoff) {
			filesToDelete = append(filesToDelete, fileStat)
		}
	}

	fmt.Printf("Found %d files to delete in %v\n", len(filesToDelete), subPath)

	for _, fileToDelete := range filesToDelete {
		err = os.RemoveAll(fileToDelete.Name())
		if err != nil {
			return err
		}
	}

	return nil
}
