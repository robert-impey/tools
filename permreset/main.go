package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/pflag"
)

func main() {
	var scriptsDir string
	pflag.StringVarP(&scriptsDir, "scriptsDirectory", "s", "", "Path to the scripts directory")
	pflag.Parse()

	if runtime.GOOS == "windows" {
		fmt.Fprintln(os.Stderr, "This program should only be run on Linux or macOs")
		os.Exit(1)
	}

	if scriptsDir == "" {
		fmt.Fprintln(os.Stderr, "The --scriptsDirectory path must be provided.")
		pflag.Usage()
		os.Exit(2)
	}

	files, err := FindFilesWithShebang(scriptsDir)
	if err != nil {
		log.Fatalf("failed to find files: %v", err)
	}

	log.Printf("Found %d files with shebangs", len(files))

	if err := ApplyPerms(files); err != nil {
		log.Printf("Completed with error: %v", err)
	}
}
