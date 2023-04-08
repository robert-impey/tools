/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"hash/maphash"
	"math/rand"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	"github.com/robert-impey/staydeleted/sdlib"
	"github.com/spf13/cobra"
)

// sweepNightlyCmd represents the sweepNightly command
var sweepNightlyCmd = &cobra.Command{
	Use:   "sweepNightly",
	Short: "Runs the stay deleted command with the nightly list",
	Long: `The list of directories to sweep can be saved in a text file.
	
A random delay can be set before starting the sweep.`,
	Run: func(cmd *cobra.Command, args []string) {
		sweepNightly()
	},
}

var Sleep int32

func init() {
	rootCmd.AddCommand(sweepNightlyCmd)
	sweepNightlyCmd.Flags().Int32VarP(&Sleep, "sleep", "s", 0,
		"The maximum number of seconds to sleep before starting. A random time during the period is chosen.")
}

func sweepNightly() {
	if Sleep > 0 {
		r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))
		wait := r.Int31n(Sleep)
		time.Sleep(time.Duration(wait) * time.Second)
	}

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
}
