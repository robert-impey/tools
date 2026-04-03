package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// GenerateRcifScript generates a PowerShell script for a single synch file
// The synch file format:
// Line 1: source path
// Line 2: destination path
// Line 3: blank
// Line 4+: files (one per line)
func GenerateRcifScript(synchFilePath, autoGenDir, scriptName string) error {
	if synchFilePath == "" {
		return fmt.Errorf("synch file path is required")
	}
	if autoGenDir == "" {
		return fmt.Errorf("autogen folder is required")
	}
	if scriptName == "" {
		return fmt.Errorf("script name is required")
	}

	if _, err := os.Stat(autoGenDir); os.IsNotExist(err) {
		return fmt.Errorf("auto-generated folder does not exist: %s", autoGenDir)
	}

	b, err := os.ReadFile(synchFilePath)
	if err != nil {
		return fmt.Errorf("failed to read synch file: %w", err)
	}

	// split into lines and normalise CRLF
	rawLines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")
	if len(rawLines) < 4 {
		return fmt.Errorf("the file must contain at least four lines: source, destination, blank, and at least one file")
	}

	id := strings.TrimSuffix(filepath.Base(synchFilePath), filepath.Ext(synchFilePath))
	source := strings.TrimSpace(rawLines[0])
	destination := strings.TrimSpace(rawLines[1])

	var files []string
	for _, line := range rawLines[3:] {
		t := strings.TrimSpace(line)
		if t == "" {
			continue
		}
		files = append(files, t)
	}

	if source == "" {
		return fmt.Errorf("the source path cannot be empty")
	}
	if destination == "" {
		return fmt.Errorf("the destination path cannot be empty")
	}
	if len(files) == 0 {
		return fmt.Errorf("at least one file must be specified")
	}

	outputScriptPath := filepath.Join(autoGenDir, fmt.Sprintf("%s.ps1", scriptName))
	if _, err := os.Stat(outputScriptPath); err == nil {
		if err := os.Remove(outputScriptPath); err != nil {
			return fmt.Errorf("failed to remove existing script: %w", err)
		}
	}

	f, err := os.Create(outputScriptPath)
	if err != nil {
		return fmt.Errorf("failed to create script file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	if _, err := fmt.Fprintln(w, "# AUTOGEN'D - DO NOT EDIT!"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "# Written %s\n\n", time.Now().Format(time.RFC1123Z)); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "param("); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "    [Parameter (Mandatory = $False)]"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "    [switch]$logged = $False"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, ")"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, `Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1"`); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}

	sourceClean := cleanFolderPathForLogName(source)
	destinationClean := cleanFolderPathForLogName(destination)

	first := true
	for _, file := range files {
		if first {
			first = false
		} else {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}

		if file == "" || strings.HasPrefix(file, "#") {
			continue
		}

		fmt.Fprintln(w, "SynchSingleFile2Ways `")
		fmt.Fprintf(w, "    -id \"%s\" `\n", id)
		fmt.Fprintf(w, "    -file \"%s\" `\n", file)
		fmt.Fprintf(w, "    -sourceFolder \"%s\" `\n", source)
		fmt.Fprintf(w, "    -sourceLogName \"%s\" `\n", sourceClean)
		fmt.Fprintf(w, "    -destinationFolder \"%s\" `\n", destination)
		fmt.Fprintf(w, "    -destinationLogName \"%s\" `\n", destinationClean)
		fmt.Fprintln(w, "    -logged $logged")
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func cleanFolderPathForLogName(path string) string {
	p := strings.TrimSpace(path)
	// normalize slashes
	p = strings.ReplaceAll(p, "/", "\\")
	// replace backslashes with underscores
	p = strings.ReplaceAll(p, "\\", "_")

	// remove invalid filename characters
	// conservative set: <>:"/\|?* and control chars
	invalid := regexp.MustCompile(`[<>:\"/\\|?*\x00-\x1F]`)
	p = invalid.ReplaceAllString(p, "")

	return p
}
