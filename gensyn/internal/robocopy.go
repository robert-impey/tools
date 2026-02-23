package internal

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	common "github.com/robert-impey/tools/internal"
)

// GenerateRobocopyScripts creates PowerShell robocopy sync scripts for all location pairs.
// The manager argument comes from the shared package; we keep the
// implementation here because the script generation logic is specific to
// the "gensyn" binary.
func GenerateRobocopyScripts(manager *common.FolderManager, autoGenFolder string) error {
	if autoGenFolder == "" {
		return fmt.Errorf("auto-generated folder cannot be empty")
	}
	if _, err := os.Stat(autoGenFolder); os.IsNotExist(err) {
		return fmt.Errorf("auto-generated folder does not exist: %s", autoGenFolder)
	}

	for _, loc1 := range manager.Locations {
		for _, loc2 := range manager.Locations {
			if loc1 == loc2 {
				continue
			}
			if err := generateScriptsForPair(manager, autoGenFolder, loc1, loc2); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateScriptsForPair(manager *common.FolderManager, autoGenFolder, src, dst string) error {
	scriptDir := getScriptPath(autoGenFolder, src, dst)

	var commonFolders []string
	for _, folder := range manager.Folders {
		if pathExists(filepath.Join(src, folder)) && pathExists(filepath.Join(dst, folder)) {
			commonFolders = append(commonFolders, folder)

			err := createRobocopySyncScript(
				folder,
				filepath.Join(scriptDir, folder+".ps1"),
				src, dst,
			)
			if err != nil {
				return err
			}
		}
	}

	if len(commonFolders) > 0 {
		return createAllFoldersRobocopySyncScript(
			commonFolders,
			filepath.Join(scriptDir, "_all.ps1"),
			src, dst,
		)
	}
	return nil
}

func createRobocopySyncScript(folder, scriptPath, src, dst string) error {
	if err := ensureParentDir(scriptPath); err != nil {
		return err
	}

	fmt.Println("Generating", scriptPath)

	f, err := os.Create(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to create script file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	if err := writePSHeader(w); err != nil {
		return err
	}

	fmt.Fprintln(w, `Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1"`)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "$folder = \"%s\"\n\n", folder)
	fmt.Fprintf(w, "$src = \"%s\"\n", filepath.Clean(src))
	fmt.Fprintf(w, "$dst = \"%s\"\n", filepath.Clean(dst))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Synch $folder $src $dst $logged")

	return w.Flush()
}

func createAllFoldersRobocopySyncScript(folders []string, scriptPath, src, dst string) error {
	if err := ensureParentDir(scriptPath); err != nil {
		return err
	}

	fmt.Println("Generating", scriptPath)

	f, err := os.Create(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to create script file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	if err := writePSHeader(w); err != nil {
		return err
	}

	fmt.Fprintln(w, `Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1"`)
	fmt.Fprintln(w)

	fmt.Fprint(w, "$folders = ")
	for i, folder := range folders {
		if i > 0 {
			fmt.Fprint(w, ", ")
		}
		fmt.Fprintf(w, "\"%s\"", folder)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w)

	absSrc, _ := filepath.Abs(src)
	absDst, _ := filepath.Abs(dst)
	fmt.Fprintf(w, "$src = \"%s\"\n", absSrc)
	fmt.Fprintf(w, "$dst = \"%s\"\n", absDst)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "foreach ($folder in $folders) {")
	fmt.Fprintln(w, "    Synch $folder $src $dst $logged")
	fmt.Fprintln(w, "}")

	return w.Flush()
}

// --- PowerShell header helpers ---

func writePSHeader(w *bufio.Writer) error {
	if err := writeAutoGenHeader(w); err != nil {
		return err
	}
	return writeSynchScriptFileParams(w)
}

func writeAutoGenHeader(w *bufio.Writer) error {
	if _, err := fmt.Fprintln(w, "# AUTOGEN'D FILE - DO NOT EDIT"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "# Created: %s\n\n", time.Now().Local().Format(time.RFC1123Z)); err != nil {
		return err
	}
	return nil
}

func writeSynchScriptFileParams(w *bufio.Writer) error {
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
	return nil
}

// --- Path helpers ---

func getScriptPath(autoGenFolder, location1, location2 string) string {
	return filepath.Join(autoGenFolder, getCleanLocationName(location1), getCleanLocationName(location2))
}

var locationCleanRegexp = regexp.MustCompile(`[:\\\/ ]+`)

func getCleanLocationName(location string) string {
	cleaned := locationCleanRegexp.ReplaceAllString(location, "_")
	return strings.TrimRight(cleaned, "_")
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ensureParentDir(path string) error {
	parent := filepath.Dir(path)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		if err := os.MkdirAll(parent, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", parent, err)
		}
	}
	return nil
}
