package cmd

/*
Copyright © 2023 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const (
	ZshPath  string = "/usr/bin/zsh"
	PwshPath string = "/snap/bin/pwsh"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print out a crontab",
	Long: `This command prints out a crontab with the daily tasks.
`,
	Run: func(cmd *cobra.Command, args []string) {
		printAllTasks()
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}

func printAllTasks() {
	printHeaderComment()

	printStayDeletedRun(0, 2)
	printSynch(2, 4)
	fmt.Println()

	printStayDeletedRun(7, 9)

	printBuild(9)
	printResetPerms(10)
	printLogsDeleter(11)
	printListManagedFolders(12)
	fmt.Println()

	printSynch(13, 4)
	fmt.Println()

	printStayDeletedRun(17, 19)

	printBuild(19)
	printResetPerms(20)
	printTidyFolder(21)
	fmt.Println()

	printStayDeletedRun(22, 24)
}

func printHeaderComment() {
	now := time.Now().UTC()
	fmt.Printf("# Cron tasks created by cron-setter at %s\n\n", now.Format("2006-01-02 15:04:05"))
}

func getLocalScripts() string {
	if envPath := os.Getenv("LOCAL_SCRIPTS"); envPath != "" {
		return envPath
	}

	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(nil)
	}
	return filepath.Join(currentUser.HomeDir, "local-scripts")
}

func getBuildScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "build", "zsh-cron-runner.sh")
}

func getListManagedFoldersScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "Write-ManagedFoldersList.ps1")
}

func getResetPermsScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "reset_perms", "reset-perms-linux.sh")
}

func getTidyFolderScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "tidy_folder", "Search-ManagedFolders.ps1")
}

func getStayDeletedScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "stay_deleted", "Sweep-ManagedFolders.ps1")
}

func getSynchScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "synch", "run-nightly.sh")
}

func getLogsDeleterScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "logs_deleter", "Clear-OldLogs.ps1")
}

func printStayDeletedRun(startHour int, endHour int) {
	script := getStayDeletedScript()

	for i := startHour; i < endHour; i++ {
		stayDeletedMinutes := rand.Int31n(60)
		fmt.Printf("%d %d * * * %s %s\n",
			stayDeletedMinutes, i, PwshPath, script)
	}
	fmt.Println()
}

func printLogsDeleter(hour int) {
	minutes := rand.Int31n(60)

	script := getLogsDeleterScript()
	fmt.Printf("%d %d * * * %s %s\n",
		minutes, hour, PwshPath, script)
}

func printResetPerms(hour int) {
	minutes := rand.Int31n(60)

	script := getResetPermsScript()
	fmt.Printf("%d %d * * * %s %s\n",
		minutes, hour, ZshPath, script)
}

func printTidyFolder(hour int) {
	minutes := rand.Int31n(60)

	script := getTidyFolderScript()
	fmt.Printf("%d %d * * * %s %s\n",
		minutes, hour, PwshPath, script)
}

func printSynch(earliestHour int32, hoursRange int32) {
	synchScript := getSynchScript()

	synchMinutes := rand.Int31n(60)
	synchHours := rand.Int31n(hoursRange) + earliestHour

	fmt.Printf("%d %d * * * %s %s\n",
		synchMinutes, synchHours, ZshPath, synchScript)
}

func printListManagedFolders(hour int) {
	minutes := rand.Int31n(60)

	script := getListManagedFoldersScript()
	fmt.Printf("%d %d * * * %s %s\n", minutes, hour, PwshPath, script)
}

func printBuild(hour int) {
	script := getBuildScript()
	minutes := rand.Int31n(60)

	fmt.Printf("%d %d * * * %s %s\n", minutes, hour, ZshPath, script)
}
