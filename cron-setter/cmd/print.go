package cmd

/*
Copyright © 2023 Robert Impey robert.impey@hotmail.co.uk
*/

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"os/user"
	"path/filepath"
	"time"
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
	printStayDeleted(dev)
	printListManagedFolders(dev)
	printResetPerms(dev)
	printLogsDeleter(dev)
	printSynch(dev)
}

func printHeaderComment() {
	now := time.Now().UTC()
	fmt.Printf("# Cron tasks created by cron-setter on %s\n\n", now.Format("2006-01-02"))
}

func printStayDeleted(dev bool) {
	runStayDeleted := getRunStayDeleted(dev)
	fmt.Println("# Stay Deleted")
	printStayDeletedRun(runStayDeleted, 0, 2)
	printStayDeletedRun(runStayDeleted, 8, 11)
	printStayDeletedRun(runStayDeleted, 19, 24)
}

func getRunStayDeleted(dev bool) string {
	executablesDir := getExecutablesDir(dev)
	return filepath.Join(executablesDir, "run-stay-deleted")
}

func printStayDeletedRun(runStayDeleted string, startHour int, endHour int) {
	for i := startHour; i < endHour; i++ {
		stayDeletedMinutes := rand.Int31n(60)
		fmt.Printf("%d %d * * * %s sweepNightly\n",
			stayDeletedMinutes, i, runStayDeleted)
	}
	fmt.Println()
}

func getExecutablesDir(dev bool) string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(nil)
	}
	buildType := getBuildType(dev)
	return filepath.Join(currentUser.HomeDir, "executables", "Linux", buildType, "x64")
}

func printResetPerms(dev bool) {
	resetPermsMinutes := rand.Int31n(60)

	flag := ""

	if dev {
		flag = " --dev"
	}

	resetPermsScript := getResetPermsScript()
	fmt.Printf("%d 11 * * * /usr/bin/zsh %s%s\n",
		resetPermsMinutes, resetPermsScript, flag)
}

func getResetPermsScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "reset-perms", "reset-perms.sh")
}

func getLocalScripts() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(nil)
	}
	return filepath.Join(currentUser.HomeDir, "local-scripts")
}

func printLogsDeleter(dev bool) {
	logsDeleterMinutes := rand.Int31n(60)

	logsDeleter := getLogsDeleter(dev)
	fmt.Printf("%d 12 * * * %s sweepAll\n",
		logsDeleterMinutes, logsDeleter)
}

func getLogsDeleter(dev bool) string {
	executablesDir := getExecutablesDir(dev)
	return filepath.Join(executablesDir, "logs-deleter")
}

func printSynch(dev bool) {
	fmt.Println()
	fmt.Println("# Synch")
	flag := ""

	if dev {
		flag = " --dev"
	}

	synchScript := getSynchScript()

	printSynchLine(2, 4, synchScript, flag)
	printSynchLine(13, 4, synchScript, flag)
}

func printSynchLine(earliestHour int32, hoursRange int32, synchScript string, flag string) {
	synchMinutes := rand.Int31n(60)
	synchHours := rand.Int31n(hoursRange) + earliestHour

	fmt.Printf("%d %d * * * /usr/bin/zsh %s%s\n",
		synchMinutes, synchHours, synchScript, flag)
}

func getSynchScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "synch", "run-nightly.sh")
}

func getBuildType(dev bool) string {
	if dev {
		return "dev"
	}
	return "prod"
}

func printListManagedFolders(dev bool) {
	minutes := rand.Int31n(60)

	exe := getListManagedFolders(dev)
	fmt.Printf("%d 6 * * * %s\n", minutes, exe)
}

func getListManagedFolders(dev bool) string {
	executablesDir := getExecutablesDir(dev)
	return filepath.Join(executablesDir, "ListManagedFolders")
}
