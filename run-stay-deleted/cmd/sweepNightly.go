/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/robert-impey/staydeleted/sdlib"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"
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
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	machineLSDir := path.Join(localScriptsDirectory, hostname)

	machineLSDirStat, err := os.Stat(machineLSDir)
	if err != nil {
		panic(err)
	}

	if !machineLSDirStat.IsDir() {
		panic(errors.New(fmt.Sprintf("%s is not a directory!", machineLSDir)))
	}
	absMachineLSDir, err := filepath.Abs(machineLSDir)
	if err != nil {
		panic(err)
	}

	userInfo, err := user.Current()

	userMachineLSDir := path.Join(absMachineLSDir, userInfo.Username)

	userMachineLSNightly := path.Join(userMachineLSDir, "staydeleted", "nightly.txt")
	machineLSNightly := path.Join(machineLSDir, "staydeleted", "nightly.txt")

	nightlyFile := ""
	nightlyErr := errors.New("No nightly file found!")

	_, err = os.Stat(userMachineLSNightly)
	if err == nil {
		nightlyFile, err = filepath.Abs(userMachineLSNightly)
		if err != nil {
			panic(err)
		}
		nightlyErr = nil

	} else {
		_, err = os.Stat(machineLSNightly)
		if err != nil {
			panic(err)
		}
		nightlyFile, err = filepath.Abs(machineLSNightly)
		if err != nil {
			panic(err)
		}
		nightlyErr = nil
	}

	if nightlyErr != nil {
		panic(nightlyErr)
	}

	fmt.Printf("Using %s\n", nightlyFile)

	const toolName = "staydeleted"
	var logsDir = filepath.Join(userInfo.HomeDir, "logs", toolName)

	if _, err := os.Stat(logsDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	timeStr := time.Now().Format("2006-01-02_15.04.05")
	outLogFileName := filepath.Join(logsDir, fmt.Sprintf("%s.log", timeStr))
	errLogFileName := filepath.Join(logsDir, fmt.Sprintf("%s.err", timeStr))

	outLogFile, err := os.Create(outLogFileName)
	if err != nil {
		panic(err)
	}
	errLogFile, err := os.Create(errLogFileName)
	if err != nil {
		panic(err)
	}

	sdlib.SweepFrom(nightlyFile, 12, outLogFile, errLogFile, false)
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

	user, err := user.Current()
	if err != nil {
		return "", err
	}
	localScriptsDir := filepath.Join(user.HomeDir, "local-scripts")
	localScriptsDirStat, err := os.Stat(localScriptsDir)
	if err != nil {
		return "", err
	}

	if !localScriptsDirStat.IsDir() {
		return "", err
	}
	absLocalScriptsDir, err := filepath.Abs(localScriptsDir)
	if err != nil {
		return "", err
	}

	return absLocalScriptsDir, nil

	return "", errors.New("Unable to find the local scripts dir!")
}
