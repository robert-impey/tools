package cmd

/*
Copyright © 2023 Robert Impey robert.impey@hotmail.co.uk
*/

import (
	"errors"
	"fmt"
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
	localScriptsDirectory, err := getLocalScriptsDirectory()
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

	userMachineLSDir := path.Join(absMachineLSDir, userInfo.Username)

	userMachineLSNightly := path.Join(userMachineLSDir, "staydeleted", "nightly.txt")
	machineLSNightly := path.Join(machineLSDir, "staydeleted", "nightly.txt")

	nightlyFile := ""
	nightlyErr := errors.New("No nightly file found!")

	_, err = os.Stat(userMachineLSNightly)
	if err == nil {
		nightlyFile, err = filepath.Abs(userMachineLSNightly)
		if err != nil {
			log.Fatalln(err)
		}
		nightlyErr = nil
	} else {
		_, err = os.Stat(machineLSNightly)
		if err == nil {
			nightlyFile, err = filepath.Abs(machineLSNightly)
			if err != nil {
				log.Fatalln(err)
			}
			nightlyErr = nil
		}
	}

	if nightlyFile == "" {
		nightlyFile, err = getManagedFoldersFile()
		if err != nil {
			log.Fatalln(err)
		}
		nightlyErr = nil
	}

	if nightlyErr != nil {
		log.Fatalln(nightlyErr)
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

func getManagedFoldersFile() (string, error) {
	theUser, err := user.Current()
	if err != nil {
		return "", err
	}

	managedFoldersFile := path.Join(theUser.HomeDir, "autogen", "managed-folders.txt")

	_, err = os.Stat(managedFoldersFile)
	if err != nil {
		return "", err
	}

	return filepath.Abs(managedFoldersFile)
}
