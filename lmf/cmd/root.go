package cmd

/*
Copyright © 2026 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/robert-impey/tools/lmf/lib"
	"github.com/spf13/cobra"
)

// Flag variables
var (
	locationsFile      string
	foldersFile        string
	managedFoldersFile string
)

var rootCmd = &cobra.Command{
	Use:   "foldermanager",
	Short: "Manages and lists directories across multiple locations",
	Long:  `A CLI tool to scan specified locations for a set of folder names and list the valid directories found.`,
	// This is the equivalent of your call() method
	RunE: func(cmd *cobra.Command, args []string) error {

		// 1. Validation Logic
		if locationsFile == "" {
			fmt.Println("Locations file is required. Use -l or --locations to specify it.")
			return nil // Return nil to avoid showing usage error if this is expected behavior
		}

		if foldersFile == "" {
			fmt.Println("Folders file is required. Use -f or --folders to specify it.")
			return nil
		}

		fmt.Printf("Reading locations from: %s\n", locationsFile)
		fmt.Printf("Reading folders from: %s\n", foldersFile)

		// 2. Load Data (Simplified helper calls)
		locations, err := readLines(locationsFile)
		if err != nil {
			return fmt.Errorf("failed to read locations: %w", err)
		}

		folders, err := readLines(foldersFile)
		if err != nil {
			return fmt.Errorf("failed to read folders: %w", err)
		}

		fm := &lib.FolderManager{
			Locations: locations,
			Folders:   folders,
		}

		// 3. Setup Output (PrintWriter equivalent)
		var outFile io.Writer = os.Stdout

		if managedFoldersFile != "" {
			fmt.Printf("Writing the list of managed folders to: %s\n", managedFoldersFile)

			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(managedFoldersFile), 0755); err != nil {
				return err
			}

			// Create/Truncate file
			f, err := os.Create(managedFoldersFile)
			if err != nil {
				return err
			}
			defer f.Close() // Automatically close when RunE finishes

			outFile = f

			// FolderManager.writeAutoGenHeader(outFile)
			writeAutoGenHeader(outFile)
		}

		// 4. Execution

		fm.ListManagedFolders(outFile)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Equivalent to picocli @Option
	rootCmd.Flags().StringVarP(&locationsFile, "locations", "l", "", "The locations file")
	rootCmd.Flags().StringVarP(&foldersFile, "folders", "f", "", "The folders file")
	rootCmd.Flags().StringVarP(&managedFoldersFile, "managed-folders-file", "m", "", "The managed folders file")
}

// Helper to read file lines
func readLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var result []string
	for _, l := range lines {
		if t := strings.TrimSpace(l); t != "" {
			result = append(result, t)
		}
	}
	return result, nil
}
func writeAutoGenHeader(outFile io.Writer) {
	// 1. Print the header line
	fmt.Fprintln(outFile, "# AUTOGEN'D FILE - DO NOT EDIT")

	// 2. Get current time in the system's local time zone
	now := time.Now()

	// 3. Format the time using the RFC1123 standard (similar to RFC_1123_DATE_TIME)
	// Example output: Mon, 02 Jan 2006 15:04:05 MST
	formattedDateTime := now.Format(time.RFC1123)

	// 4. Print the formatted string with an extra newline at the end
	fmt.Fprintf(outFile, "# Created: %s\n\n", formattedDateTime)
}
