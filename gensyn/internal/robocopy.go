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

	fmt.Fprintf(w, "$folder = \"%s\"\n\n", folder)
	fmt.Fprintf(w, "$src = \"%s\"\n", filepath.Clean(src))
	fmt.Fprintf(w, "$dst = \"%s\"\n", filepath.Clean(dst))
	fmt.Fprintln(w)

	writeSynchBody(w)

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
	writeSynchBody(w)
	fmt.Fprintln(w, "}")

	return w.Flush()
}

func writeSynchBody(w *bufio.Writer) {
	fmt.Fprintln(w, "$srcLogStr = $src -replace '[:\\\\/ ]+', '_'")
	fmt.Fprintln(w, "$dstLogStr = $dst -replace '[:\\\\/ ]+', '_'")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "$srcPath = \"$($src)\\$($folder)\"")
	fmt.Fprintln(w, "$dstPath = \"$($dst)\\$($folder)\"")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "if ((Test-Path $srcPath) -and (Test-Path $dstPath))")
	fmt.Fprintln(w, "{")
	fmt.Fprintln(w, "    Write-Output \"$($srcPath) and $($dstPath) exist - synching...\"")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "    if ($logged)")
	fmt.Fprintln(w, "    {")
	fmt.Fprintln(w, "        $logTimeStr = Get-Date -Format \"yyyy-MM-ddTHH_mm_ss\"")
	fmt.Fprintln(w, "        $logFileBase = \"$($logTimeStr)-$($srcLogStr)-$($dstLogStr)-$($folder).robocopy-synch\"")
	fmt.Fprintln(w, "        $logPathBase = \"$($env:USERPROFILE)\\logs\\synch\\$($logFileBase)\"")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "        $logFile = \"$($logPathBase).log\"")
	fmt.Fprintln(w, "        $errFile = \"$($logPathBase).err\"")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "        Start-Process ROBOCOPY -ArgumentList \"\"\"$($srcPath)\"\" \"\"$($dstPath)\"\" /E /XO /R:0\" `")
	fmt.Fprintln(w, "            -RedirectStandardOutput $logFile `")
	fmt.Fprintln(w, "            -RedirectStandardError $errFile `")
	fmt.Fprintln(w, "            -NoNewWindow `")
	fmt.Fprintln(w, "            -Wait")
	fmt.Fprintln(w, "    } else")
	fmt.Fprintln(w, "    {")
	fmt.Fprintln(w, "        Start-Process ROBOCOPY -ArgumentList \"\"\"$($srcPath)\"\" \"\"$($dstPath)\"\" /E /XO /R:0\" `")
	fmt.Fprintln(w, "            -NoNewWindow `")
	fmt.Fprintln(w, "            -Wait")
	fmt.Fprintln(w, "    }")
	fmt.Fprintln(w, "} else")
	fmt.Fprintln(w, "{")
	fmt.Fprintln(w, "    Write-Output \"Check if $($srcPath) and $($dstPath) exist\"")
	fmt.Fprintln(w, "}")
}

// --- PowerShell header helpers ---

func writePSHeader(w *bufio.Writer) error {
	if err := common.WriteHeader(w, "# AUTOGEN'D FILE - DO NOT EDIT"); err != nil {
		return err
	}
	return writeSynchScriptFileParams(w)
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
