package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

func main() {
	var scriptsDir string
	flag.StringVar(&scriptsDir, "s", "", "Path to the scripts directory")
	flag.StringVar(&scriptsDir, "scriptsDirectory", "", "Path to the scripts directory")
	flag.Parse()

	if runtime.GOOS == "windows" {
		fmt.Fprintln(os.Stderr, "This program should only be run on Linux or macOs")
		os.Exit(1)
	}

	if scriptsDir == "" {
		fmt.Fprintln(os.Stderr, "The --scriptsDirectory path must be provided.")
		flag.Usage()
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
