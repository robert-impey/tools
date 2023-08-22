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
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dev, _ := cmd.Flags().GetBool("dev")
		print(dev)
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	printCmd.Flags().BoolP("dev", "d", false, "add dev versions to cron")
}

func print(dev bool) {
	printStayDeleted(dev)
	printResetPerms(dev)
	printLogsDeleter(dev)
	printSynch(dev)
}

func printStayDeleted(dev bool) {
	runStayDeleted := getRunStayDeleted(dev)
	fmt.Println("# Stay Deleted")
	printStayDeletedRun(runStayDeleted, 0, 2)
	printStayDeletedRun(runStayDeleted, 6, 11)
	printStayDeletedRun(runStayDeleted, 19, 23)
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
	user, err := user.Current()
	if err != nil {
		log.Fatal(nil)
	}
	buildType := getBuildType(dev)
	return filepath.Join(user.HomeDir, "executables", "Linux", buildType, "x64")
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
	user, err := user.Current()
	if err != nil {
		log.Fatal(nil)
	}
	return filepath.Join(user.HomeDir, "local-scripts")
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
