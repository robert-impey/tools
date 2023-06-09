package cmd

/*
Copyright © 2023 Robert Impey robert.impey@hotmail.co.uk
*/

import (
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
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
		print()
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}

func print() {
	printStayDeleted()
	printResetPerms()
	printLogsDeleter()
	printSynch()
}

func printStayDeleted() {
	fmt.Println("# Stay Deleted")
	for i := 0; i < 11; i++ {
		stayDeletedMinutes := rand.Int31n(60)
		fmt.Printf("%d %d * * * /home/robert/executables/Linux/prod/x64/run-stay-deleted sweepNightly\n", stayDeletedMinutes, i)
	}
	for i := 18; i < 23; i++ {
		stayDeletedMinutes := rand.Int31n(60)
		fmt.Printf("%d %d * * * /home/robert/executables/Linux/prod/x64/run-stay-deleted sweepNightly\n", stayDeletedMinutes, i)
	}
	fmt.Println()
}

func printResetPerms() {
	resetPermsMinutes := rand.Int31n(60)

	fmt.Printf("%d 11 * * * /usr/bin/zsh /home/robert/local-scripts/_Common/reset-perms/reset-perms.sh\n", resetPermsMinutes)
}

func printLogsDeleter() {
	logsDeleterMinutes := rand.Int31n(60)

	fmt.Printf("%d 12 * * * /home/robert/executables/Linux/prod/x64/logs-deleter sweepAll\n", logsDeleterMinutes)
}

func printSynch() {
	synchMinutes := rand.Int31n(60)
	synchHours := rand.Int31n(4) + 13

	fmt.Printf("%d %d * * * /usr/bin/zsh /home/robert/local-scripts/_Common/synch/run-nightly.sh\n", synchMinutes, synchHours)
}
