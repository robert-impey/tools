package lib

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

const filesCmdLineTemplate = "%v %v/%v %v/"
const dirsCmdLineTemplate = "%v %v/%v/ %v/%v"

type ScriptsInfo struct {
	name, dir, synch, src, dst string
	items                      []string
}

func GenerateSynchScripts(files bool, autoGenDir string, gssFile string) error {
	fmt.Printf("Generating synch scripts for %v\n", gssFile)

	scriptInfo, err := ParseGSSFile(gssFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse %v - %v\n", gssFile, err)
		return err
	}

	if err := writeScripts(files, autoGenDir, scriptInfo); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write all items script for %v - %v\n", gssFile, err)
		return err
	}
	return nil
}

func ParseGSSFile(gssFileName string) (*ScriptsInfo, error) {
	scriptsInfo := new(ScriptsInfo)

	base := filepath.Base(gssFileName)

	scriptsInfo.name = strings.TrimSuffix(base, path.Ext(base))

	// Script directory
	dir := filepath.Dir(gssFileName)
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find the absolute directory name for %v\n", dir)
		return scriptsInfo, err
	}
	scriptsInfo.dir = absDir

	// Read the file
	gssFile, err := os.Open(gssFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %v - %v\n", gssFileName, err)
		return nil, err
	}
	defer gssFile.Close()
	input := bufio.NewScanner(gssFile)

	// Get the script root
	input.Scan()
	scriptsInfo.synch = input.Text()

	// Get the src and the dst
	input.Scan()
	scriptsInfo.src = input.Text()
	input.Scan()
	scriptsInfo.dst = input.Text()

	// Skip the blank line
	input.Scan()

	// Get the list of directories
	dirs := mapset.NewSet[string]()
	for input.Scan() {
		dir := strings.Trim(input.Text(), " /")
		if len(dir) > 0 {
			dirs.Add(dir)
		}
	}
	dirsSlice := dirs.ToSlice()
	sort.Strings(dirsSlice)
	scriptsInfo.items = dirsSlice

	return scriptsInfo, nil
}

func writeScripts(files bool, autoGenDir string, scriptsInfo *ScriptsInfo) error {
	fmt.Printf("Generating scripts in %v\n", autoGenDir)
	fmt.Printf("Synch root: %v\n", scriptsInfo.synch)
	fmt.Printf("Source: %v\n", scriptsInfo.src)
	fmt.Printf("Destination: %v\n", scriptsInfo.dst)
	fmt.Println("Directories to synch:")
	for _, dir := range scriptsInfo.items {
		fmt.Println(dir)
	}
	fmt.Println()

	scriptName := fmt.Sprintf("%s.sh", scriptsInfo.name)
	scriptFileName := filepath.Join(autoGenDir, scriptName)

	if _, err := os.Stat(scriptFileName); err == nil {
		fmt.Printf("%v exists - Deleting...\n", scriptFileName)
		err := os.Remove(scriptFileName)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
	}

	var allScriptsBuffer bytes.Buffer
	allScriptsBuffer.WriteString("#!/bin/bash\n# AUTOGEN'D - DO NOT EDIT!\n")

	allScriptsBuffer.WriteString(fmt.Sprintf("# Generated on %s\n\n", getNowFmt()))
	allScriptsBuffer.WriteString("date\n\n")

	for _, dir := range scriptsInfo.items {
		to := getCmdLine(
			files,
			scriptsInfo.synch,
			dir,
			scriptsInfo.src,
			scriptsInfo.dst)
		allScriptsBuffer.WriteString(getEchoLine(to) + "\n")
		allScriptsBuffer.WriteString(to + "\n")

		from := getCmdLine(
			files,
			scriptsInfo.synch,
			dir,
			scriptsInfo.dst,
			scriptsInfo.src)
		allScriptsBuffer.WriteString(getEchoLine(from) + "\n")
		allScriptsBuffer.WriteString(from + "\n")

		allScriptsBuffer.WriteString("\n")
	}

	allScriptsBuffer.WriteString("\ndate\n")

	err := os.WriteFile(scriptFileName, allScriptsBuffer.Bytes(), 0x755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write script to %v - %v\n", scriptFileName, err)
	}

	if files {
		return nil
	}

	if len(scriptsInfo.items) > 1 {
		fileDir := filepath.Join(autoGenDir, scriptsInfo.name)

		if _, err := os.Stat(fileDir); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(fileDir, os.ModePerm)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				return err
			}
		}

		for _, dir := range scriptsInfo.items {
			scriptName := fmt.Sprintf("%s.sh", dir)
			scriptFileName := filepath.Join(fileDir, scriptName)

			if _, err := os.Stat(scriptFileName); err == nil {
				fmt.Printf("%v exists - Deleting...\n", scriptFileName)
				err := os.Remove(scriptFileName)
				if err != nil {
					fmt.Fprint(os.Stderr, err.Error())
				}
			}

			var scriptBuffer bytes.Buffer
			scriptBuffer.WriteString("#!/bin/bash\n# AUTOGEN'D - DO NOT EDIT!\n")
			scriptBuffer.WriteString(fmt.Sprintf("# Generated on %s\n\n", getNowFmt()))

			scriptBuffer.WriteString("date\n\n")

			to := getCmdLine(
				false,
				scriptsInfo.synch,
				dir,
				scriptsInfo.src,
				scriptsInfo.dst)
			scriptBuffer.WriteString(getEchoLine(to) + "\n")
			scriptBuffer.WriteString(to + "\n")

			from := getCmdLine(
				false,
				scriptsInfo.synch,
				dir,
				scriptsInfo.dst,
				scriptsInfo.src)
			scriptBuffer.WriteString(getEchoLine(from) + "\n")
			scriptBuffer.WriteString(from + "\n")

			scriptBuffer.WriteString("\n")

			scriptBuffer.WriteString("\ndate\n")

			err = os.WriteFile(scriptFileName, scriptBuffer.Bytes(), 0x755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to write script to %v - %v\n", scriptFileName, err)
			}
		}
	}

	return nil
}

func getCmdLine(files bool, synchRoot, item, src, dst string) string {
	if files {
		return fmt.Sprintf(
			filesCmdLineTemplate,
			synchRoot,
			src, item,
			dst)
	}
	return fmt.Sprintf(
		dirsCmdLineTemplate,
		synchRoot,
		src, item,
		dst, item)
}

func getEchoLine(cmd string) string {
	return fmt.Sprintf("echo '%s'", cmd)
}

func getNowFmt() string {
	return time.Now().UTC().Format(time.RFC1123)
}
