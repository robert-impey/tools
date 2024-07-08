package cmd

/*
Copyright © 2023 Robert Impey robert.impey@hotmail.co.uk
*/

import (
	"errors"
	"fmt"
	"github.com/robert-impey/tools/managed-folders/mflib"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	"github.com/robert-impey/tools/staydeleted/sdlib"
	"github.com/spf13/cobra"
)

// sweepNightlyCmd represents the sweepNightly command
var sweepNightlyCmd = &cobra.Command{
	Use:   "sweepNightly",
	Short: "Runs the stay deleted command with the nightly list",
	Long: `The list of directories to sweep can be saved in a text file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		find, err := cmd.PersistentFlags().GetBool("Find")
		if err != nil {
			log.Fatalln(err)
		}
		err = sweepNightly(find)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sweepNightlyCmd)

	sweepNightlyCmd.PersistentFlags().BoolP("Find", "f", false,
		"Just find the file to sweep from and quit")
}

func sweepNightly(find bool) error {
	localScriptsDirectory, err := mflib.GetLocalScriptsDirectory()
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	machineLSDir := path.Join(localScriptsDirectory, hostname)

	machineLSDirStat, err := os.Stat(machineLSDir)
	if err != nil {
		return err
	}

	if !machineLSDirStat.IsDir() {
		log.Fatalf("%s is not a directory!\n", machineLSDir)
	}
	absMachineLSDir, err := filepath.Abs(machineLSDir)
	if err != nil {
		log.Fatalln(err)
	}

	userInfo, err := user.Current()
	if err != nil {
		return err
	}

	userMachineLSDir := path.Join(absMachineLSDir, userInfo.Username)

	userMachineLSNightly := path.Join(userMachineLSDir, "staydeleted", "nightly.txt")
	machineLSNightly := path.Join(machineLSDir, "staydeleted", "nightly.txt")

	nightlyFile := ""

	_, err = os.Stat(userMachineLSNightly)
	if err == nil {
		nightlyFile, err = filepath.Abs(userMachineLSNightly)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		_, err = os.Stat(machineLSNightly)
		if err == nil {
			nightlyFile, err = filepath.Abs(machineLSNightly)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	if nightlyFile == "" {
		nightlyFile, err = mflib.GetManagedFoldersFileName()
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("Using %s\n", nightlyFile)

	if find {
		return nil
	}

	const toolName = "staydeleted"
	var logsDir = filepath.Join(userInfo.HomeDir, "logs", toolName)

	if _, err := os.Stat(logsDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(logsDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	timeStr := time.Now().Format("2006-01-02_15.04.05")
	outLogFileName := filepath.Join(logsDir, fmt.Sprintf("%s.log", timeStr))
	errLogFileName := filepath.Join(logsDir, fmt.Sprintf("%s.err", timeStr))

	outLogFile, err := os.Create(outLogFileName)
	if err != nil {
		return err
	}
	errLogFile, err := os.Create(errLogFileName)
	if err != nil {
		return err
	}

	err = sdlib.SweepFrom(nightlyFile, 12, outLogFile, errLogFile, false)

	return err
}
