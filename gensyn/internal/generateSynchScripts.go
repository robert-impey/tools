package internal

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

func writeScriptFile(path string, content string) error {
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(content), 0o755)
}

func buildScriptHeader() string {
	return fmt.Sprintf(
		"#!/bin/bash\n# AUTOGEN'D - DO NOT EDIT!\n# Generated on %s\n\ndate\n\n",
		getNowFmt(),
	)
}

func buildScriptBody(files bool, info *ScriptsInfo, item string) string {
	var buf strings.Builder

	to := getCmdLine(files, info.synch, item, info.src, info.dst)
	from := getCmdLine(files, info.synch, item, info.dst, info.src)

	buf.WriteString(getEchoLine(to) + "\n")
	buf.WriteString(to + "\n")
	buf.WriteString(getEchoLine(from) + "\n")
	buf.WriteString(from + "\n\n")

	return buf.String()
}

func buildFullScript(files bool, info *ScriptsInfo, item string) string {
	var buf strings.Builder
	buf.WriteString(buildScriptHeader())
	buf.WriteString(buildScriptBody(files, info, item))
	buf.WriteString("date\n")
	return buf.String()
}

func writeScripts(files bool, autoGenDir string, info *ScriptsInfo) error {
	fmt.Printf("Generating scripts in %v\n", autoGenDir)
	fmt.Printf("Synch root: %v\n", info.synch)
	fmt.Printf("Source: %v\n", info.src)
	fmt.Printf("Destination: %v\n", info.dst)

	// Print items
	label := "Directories"
	if files {
		label = "Files"
	}
	fmt.Printf("%s to synch:\n", label)
	for _, item := range info.items {
		fmt.Println(item)
	}
	fmt.Println()

	// Main script
	mainScriptPath := filepath.Join(autoGenDir, info.name+".sh")
	if err := writeScriptFile(mainScriptPath, buildFullScript(files, info, "")); err != nil {
		return err
	}

	// If files mode, we're done
	if files {
		return nil
	}

	// Per-item scripts
	if len(info.items) > 1 {
		itemDir := filepath.Join(autoGenDir, info.name)
		if err := os.MkdirAll(itemDir, os.ModePerm); err != nil {
			return err
		}

		for _, item := range info.items {
			scriptPath := filepath.Join(itemDir, item+".sh")
			content := buildFullScript(false, info, item)
			if err := writeScriptFile(scriptPath, content); err != nil {
				return err
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

type FolderManager struct {
	Locations []string
	Folders   []string
}

func CreateFolderManager(locationsPath, foldersPath string) (*FolderManager, error) {
	locations, err := readLinesFromPath(locationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read locations: %w", err)
	}
	if len(locations) == 0 {
		return nil, fmt.Errorf("locations cannot be empty")
	}

	folders, err := readLinesFromPath(foldersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read folders: %w", err)
	}
	if len(folders) == 0 {
		return nil, fmt.Errorf("folders cannot be empty")
	}

	return &FolderManager{
		Locations: locations,
		Folders:   folders,
	}, nil
}

func readLinesFromPath(path string) ([]string, error) {
	// Check existence first
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", path)
		}
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func (fm *FolderManager) GenerateRobocopyScripts(autoGenFolder string) error {
	// Validate autogen folder
	if autoGenFolder == "" {
		return fmt.Errorf("auto-generated folder cannot be empty")
	}
	if _, err := os.Stat(autoGenFolder); os.IsNotExist(err) {
		return fmt.Errorf("auto-generated folder does not exist: %s", autoGenFolder)
	}

	// Iterate over all pairs of locations
	for _, location1 := range fm.Locations {
		for _, location2 := range fm.Locations {

			if location1 == location2 {
				continue
			}

			scriptPath := getScriptPath(autoGenFolder, location1, location2)
			sourcePath := location1
			destinationPath := location2

			var commonFolders []string

			for _, folder := range fm.Folders {
				src := filepath.Join(sourcePath, folder)
				dst := filepath.Join(destinationPath, folder)

				srcExists := pathExists(src)
				dstExists := pathExists(dst)

				if srcExists && dstExists {
					commonFolders = append(commonFolders, folder)

					err := createRobocopySyncScript(
						folder,
						filepath.Join(scriptPath, folder+".ps1"),
						sourcePath,
						destinationPath,
					)
					if err != nil {
						return err
					}
				}
			}

			if len(commonFolders) > 0 {
				err := createAllFoldersRobocopySyncScript(
					commonFolders,
					filepath.Join(scriptPath, "_all.ps1"),
					sourcePath,
					destinationPath,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func getScriptPath(autoGenFolder, location1, location2 string) string {
	clean1 := getCleanLocationName(location1)
	clean2 := getCleanLocationName(location2)

	return filepath.Join(autoGenFolder, clean1, clean2)
}

var locationCleanRegexp = regexp.MustCompile(`[:\\\/ ]+`)

func getCleanLocationName(location string) string {
	// Replace sequences of [: \ / space] with "_"
	cleaned := locationCleanRegexp.ReplaceAllString(location, "_")

	// Trim trailing underscores
	cleaned = strings.TrimRight(cleaned, "_")

	return cleaned
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func createRobocopySyncScript(
	folder string,
	scriptPath string,
	sourcePath string,
	destinationPath string,
) error {
	// Ensure parent directory exists
	parent := filepath.Dir(scriptPath)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		if err := os.MkdirAll(parent, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", parent, err)
		}
	}

	fmt.Println("Generating", scriptPath)

	f, err := os.Create(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to create script file: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	// Your helper functions (to be implemented)
	if err := writeAutoGenHeader(writer); err != nil {
		return err
	}
	if err := writeSynchScriptFileParams(writer); err != nil {
		return err
	}

	// PowerShell module import
	fmt.Fprintln(writer, `Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1"`)
	fmt.Fprintln(writer)

	// Variables
	fmt.Fprintf(writer, "$folder = \"%s\"\n\n", folder)
	fmt.Fprintf(writer, "$src = \"%s\"\n", filepath.Clean(sourcePath))
	fmt.Fprintf(writer, "$dst = \"%s\"\n", filepath.Clean(destinationPath))
	fmt.Fprintln(writer)

	// Synch call
	fmt.Fprintln(writer, "Synch $folder $src $dst $logged")

	return writer.Flush()
}

func writeAutoGenHeader(w *bufio.Writer) error {
	// Header line
	if _, err := fmt.Fprintln(w, "# AUTOGEN'D FILE - DO NOT EDIT"); err != nil {
		return err
	}

	// RFC 1123 timestamp (same as Java's DateTimeFormatter.RFC_1123_DATE_TIME)
	now := time.Now().Local()
	formatted := now.Format(time.RFC1123Z)

	if _, err := fmt.Fprintf(w, "# Created: %s\n\n", formatted); err != nil {
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
	if _, err := fmt.Fprintln(w); err != nil { // blank line
		return err
	}
	return nil
}

func createAllFoldersRobocopySyncScript(
	commonFolders []string,
	scriptPath string,
	sourcePath string,
	destinationPath string,
) error {
	// Ensure parent directory exists
	parent := filepath.Dir(scriptPath)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", parent, err)
	}

	fmt.Println("Generating", scriptPath)

	f, err := os.Create(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to create script file: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	// Header + params
	if err := writeAutoGenHeader(writer); err != nil {
		return err
	}
	if err := writeSynchScriptFileParams(writer); err != nil {
		return err
	}

	// PowerShell module import
	fmt.Fprintln(writer, `Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1"`)
	fmt.Fprintln(writer)

	// $folders = "A", "B", "C"
	fmt.Fprint(writer, "$folders = ")

	for i, folder := range commonFolders {
		if i > 0 {
			fmt.Fprint(writer, ", ")
		}
		fmt.Fprintf(writer, "\"%s\"", folder)
	}

	fmt.Fprintln(writer)
	fmt.Fprintln(writer)

	// $src and $dst
	absSrc, _ := filepath.Abs(sourcePath)
	absDst, _ := filepath.Abs(destinationPath)

	fmt.Fprintf(writer, "$src = \"%s\"\n", absSrc)
	fmt.Fprintf(writer, "$dst = \"%s\"\n", absDst)
	fmt.Fprintln(writer)

	// foreach loop
	fmt.Fprintln(writer, "foreach ($folder in $folders) {")
	fmt.Fprintln(writer, "    Synch $folder $src $dst $logged")
	fmt.Fprintln(writer, "}")

	return writer.Flush()
}
