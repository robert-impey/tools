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
	printBuild(dev)
	printResetPerms(dev)
	printLogsDeleter(dev)
	printSynch(dev)
	printCronSetter(dev)
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

func getFlag(dev bool) string {
	if dev {
		return " --dev"
	}

	return ""
}

func getBuildScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "build", "Build.ps1")
}

func getListManagedFoldersScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "List-ManagedFolders.ps1")
}

func getResetPermsScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "reset-perms", "reset-perms.sh")
}

func getSynchScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "synch", "run-nightly.sh")
}

func printHeaderComment() {
	now := time.Now().UTC()
	fmt.Printf("# Cron tasks created by cron-setter on %s\n\n", now.Format("2006-01-02"))
}

func printStayDeleted(dev bool) {
	runStayDeleted := getExecutable("run-stay-deleted", dev)
	fmt.Println("# Stay Deleted")
	printStayDeletedRun(runStayDeleted, 0, 2)
	printStayDeletedRun(runStayDeleted, 8, 11)
	printStayDeletedRun(runStayDeleted, 19, 23)
}

func printStayDeletedRun(runStayDeleted string, startHour int, endHour int) {
	for i := startHour; i < endHour; i++ {
		stayDeletedMinutes := rand.Int31n(60)
		fmt.Printf("%d %d * * * %s sweepNightly\n",
			stayDeletedMinutes, i, runStayDeleted)
	}
	fmt.Println()
}

func printLogsDeleter(dev bool) {
	logsDeleterMinutes := rand.Int31n(60)

	exe := getExecutable("logs-deleter", dev)
	fmt.Printf("%d 12 * * * %s sweepAll\n",
		logsDeleterMinutes, exe)
}

func printCronSetter(dev bool) {
	minutes := rand.Int31n(60)

	exe := getExecutable("cron-setter", dev)
	fmt.Printf("%d 23 * * * %s%s | /usr/bin/crontab -\n",
		minutes, exe, getFlag(dev))
}

func printResetPerms(dev bool) {
	resetPermsMinutes := rand.Int31n(60)

	resetPermsScript := getResetPermsScript()
	fmt.Printf("%d 11 * * * /usr/bin/zsh %s%s\n",
		resetPermsMinutes, resetPermsScript, getFlag(dev))
}

func printSynch(dev bool) {
	fmt.Println()
	fmt.Println("# Synch")

	flag := getFlag(dev)
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

func printListManagedFolders(dev bool) {
	script := getListManagedFoldersScript()
	minutes := rand.Int31n(60)

	fmt.Printf("%d 6 * * * /snap/bin/pwsh %s%s\n", minutes, script, getFlag(dev))
}

func printBuild(dev bool) {
	script := getBuildScript()
	minutes := rand.Int31n(60)

	fmt.Printf("%d 7 * * * /snap/bin/pwsh %s%s\n", minutes, script, getFlag(dev))
}
