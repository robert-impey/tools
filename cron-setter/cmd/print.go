package cmd

/*
Copyright © 2023 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const (
	ZshPath             string = "/usr/bin/zsh"
	PwshPath            string = "/snap/bin/pwsh"
	headerCommentPrefix string = "# Start of crontab tasks created by cron-setter"
	footerCommentPrefix string = "# End of cron tasks created by cron-setter"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print out a crontab",
	Long: `This command prints out a crontab with the daily tasks.
`,
	Run: func(cmd *cobra.Command, args []string) {
		printAllTasks(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}

func printAllTasks(w io.Writer) {
	creationTime := time.Now().UTC()
	creationTimeStr := creationTime.Format("2006-01-02 15:04:05")

	printHeaderComment(w, creationTimeStr)

	printStayDeletedRun(w, 0, 2)
	printSynch(w, 2, 4)
	printResetPerms(w, 7)
	fmt.Fprintln(w)

	printStayDeletedRun(w, 8, 10)

	printBuild(w, 10)
	printLogsDeleter(w, 11)
	printListManagedFolders(w, 12)
	fmt.Fprintln(w)

	printSynch(w, 13, 4)
	printResetPerms(w, 17)
	fmt.Fprintln(w)

	printStayDeletedRun(w, 18, 20)

	printBuild(w, 20)
	fmt.Fprintln(w)

	printStayDeletedRun(w, 21, 24)

	printFooterComment(w, creationTimeStr)
}

func printHeaderComment(w io.Writer, creationTimeStr string) {
	fmt.Fprintf(w, "%s at %s\n\n", headerCommentPrefix, creationTimeStr)
}

func printFooterComment(w io.Writer, creationTimeStr string) {
	fmt.Fprintf(w, "%s at %s\n", footerCommentPrefix, creationTimeStr)
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
	return filepath.Join(localScriptsDir, "_Common", "build", "Build.ps1")
}

func getListManagedFoldersScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "Write-ManagedFoldersList.ps1")
}

func getResetPermsScript() string {
	localScriptsDir := getLocalScripts()
	return filepath.Join(localScriptsDir, "_Common", "reset_perms", "reset-perms-linux.sh")
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

func printStayDeletedRun(w io.Writer, startHour int, endHour int) {
	script := getStayDeletedScript()

	for i := startHour; i < endHour; i++ {
		stayDeletedMinutes := rand.Int31n(60)
		fmt.Fprintf(w, "%d %d * * * %s %s\n",
			stayDeletedMinutes, i, PwshPath, script)
	}
	fmt.Fprintln(w)
}

func printLogsDeleter(w io.Writer, hour int) {
	minutes := rand.Int31n(60)

	script := getLogsDeleterScript()
	fmt.Fprintf(w, "%d %d * * * %s %s\n",
		minutes, hour, PwshPath, script)
}

func printResetPerms(w io.Writer, hour int) {
	minutes := rand.Int31n(60)

	script := getResetPermsScript()
	fmt.Fprintf(w, "%d %d * * * %s %s\n",
		minutes, hour, ZshPath, script)
}

func printSynch(w io.Writer, earliestHour int32, hoursRange int32) {
	synchScript := getSynchScript()

	synchMinutes := rand.Int31n(60)
	synchHours := rand.Int31n(hoursRange) + earliestHour

	fmt.Fprintf(w, "%d %d * * * %s %s\n",
		synchMinutes, synchHours, ZshPath, synchScript)
}

func printListManagedFolders(w io.Writer, hour int) {
	minutes := rand.Int31n(60)

	script := getListManagedFoldersScript()
	fmt.Fprintf(w, "%d %d * * * %s %s\n", minutes, hour, PwshPath, script)
}

func printBuild(w io.Writer, hour int) {
	script := getBuildScript()
	minutes := rand.Int31n(60)

	fmt.Fprintf(w, "%d %d * * * %s %s\n", minutes, hour, PwshPath, script)
}
