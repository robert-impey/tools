package cmd

/*
Copyright © 2023 Robert Impey robert.impey@hotmail.co.uk
*/

import (
	"fmt"
	"log"
	"math/rand"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print out a crontab",
	Long: `This command prints out a crontab with the daily tasks.

Optionally, this can use the development versions of the commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		dev, _ := cmd.Flags().GetBool("dev")
		printAllTasks(dev)
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	printCmd.Flags().BoolP("dev", "d", false, "add dev versions to cron")
}

func printAllTasks(dev bool) {
	printHeaderComment()

	printStayDeletedRun(0, 2, dev)
	printSynch(2, 4, dev)
	fmt.Println()

	printStayDeletedRun(7, 9, dev)

	printBuild(9)
	printResetPerms(10, dev)
	printLogsDeleter(11, dev)
	printListManagedFolders(12)
	fmt.Println()

	printSynch(13, 4, dev)
	fmt.Println()

	printStayDeletedRun(17, 19, dev)

	printBuild(19)
	printResetPerms(20, dev)
	fmt.Println()

	printStayDeletedRun(21, 24, dev)
}

func printHeaderComment() {
	now := time.Now().UTC()
	fmt.Printf("# Cron tasks created by cron-setter at %s\n\n", now.Format("2006-01-02 15:04:05"))
}

func getExecutable(name string, dev bool) string {
	executablesDir := getExecutablesDir(dev)
	return filepath.Join(executablesDir, name)
}

func getExecutablesDir(dev bool) string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(nil)
	}
	buildType := getBuildType(dev)
	return filepath.Join(currentUser.HomeDir, "executables", "Linux", buildType, "x64")
}

func getLocalScripts() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(nil)
	}
	return filepath.Join(currentUser.HomeDir, "local-scripts")
}

func getBuildType(dev bool) string {
	if dev {
		return "dev"
	}
	return "prod"
}

func getFlag(dev bool, singleDash bool) string {
	if dev {
		if singleDash {
			return " -dev"
		}
		return " --dev"
	}

	return ""
}

func getBuildScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "build", "zsh-cron-runner.sh")
}

func getListManagedFoldersScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "list-managed-folders.sh")
}

func getResetPermsScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "reset_perms", "reset-perms.sh")
}

func getSynchScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "synch", "run-nightly.sh")
}

func printStayDeletedRun(startHour int, endHour int, dev bool) {
	runStayDeleted := getExecutable("run-stay-deleted", dev)
	for i := startHour; i < endHour; i++ {
		stayDeletedMinutes := rand.Int31n(60)
		fmt.Printf("%d %d * * * %s sweepNightly\n",
			stayDeletedMinutes, i, runStayDeleted)
	}
	fmt.Println()
}

func printLogsDeleter(hour int, dev bool) {
	logsDeleterMinutes := rand.Int31n(60)

	exe := getExecutable("logs-deleter", dev)
	fmt.Printf("%d %d * * * %s sweepAll\n",
		logsDeleterMinutes, hour, exe)
}

func printResetPerms(hour int, dev bool) {
	resetPermsMinutes := rand.Int31n(60)

	resetPermsScript := getResetPermsScript()
	fmt.Printf("%d %d * * * /usr/bin/zsh %s%s\n",
		resetPermsMinutes, hour, resetPermsScript, getFlag(dev, false))
}

func printSynch(earliestHour int32, hoursRange int32, dev bool) {
	flag := getFlag(dev, false)
	synchScript := getSynchScript()

	synchMinutes := rand.Int31n(60)
	synchHours := rand.Int31n(hoursRange) + earliestHour

	fmt.Printf("%d %d * * * /usr/bin/zsh %s%s\n",
		synchMinutes, synchHours, synchScript, flag)
}

func printListManagedFolders(hour int) {
	minutes := rand.Int31n(60)

	script := getListManagedFoldersScript()
	fmt.Printf("%d %d * * * /usr/bin/zsh %s\n", minutes, hour, script)
}

func printBuild(hour int) {
	script := getBuildScript()
	minutes := rand.Int31n(60)

	fmt.Printf("%d %d * * * /usr/bin/zsh %s\n", minutes, hour, script)
}
