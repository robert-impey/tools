package internal

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

// ScriptsInfo holds parsed data from a .gss file.
type ScriptsInfo struct {
	name, dir, synch, src, dst string
	items                      []string
}

// RsyncScriptGenerator holds the configuration for generating rsync shell scripts.
type RsyncScriptGenerator struct {
	Files      bool   // true = file mode, false = directory mode
	AutoGenDir string // output directory for generated scripts
	Powershell bool   // true = generate PowerShell scripts (.ps1) instead of bash
}

// GenerateSynchScripts is a convenience wrapper preserving the original API.
func GenerateSynchScripts(files bool, autoGenDir string, gssFile string, powershell bool) error {
	g := &RsyncScriptGenerator{Files: files, AutoGenDir: autoGenDir, Powershell: powershell}
	return g.GenerateSynchScripts(gssFile)
}

// GenerateSynchScripts parses a .gss file and writes rsync shell scripts.
func (g *RsyncScriptGenerator) GenerateSynchScripts(gssFile string) error {
	fmt.Printf("Generating synch scripts for %v\n", gssFile)

	info, err := ParseGSSFile(gssFile)
	if err != nil {
		return fmt.Errorf("unable to parse %s: %w", gssFile, err)
	}

	if err := g.writeScripts(info); err != nil {
		return fmt.Errorf("unable to write scripts for %s: %w", gssFile, err)
	}
	return nil
}

// ParseGSSFile reads a .gss config file and returns structured ScriptsInfo.
func ParseGSSFile(gssFileName string) (*ScriptsInfo, error) {
	info := new(ScriptsInfo)
	if err := populateInfoFromPath(gssFileName, info); err != nil {
		return info, err
	}

	synch, src, dst, items, err := readGSSConfigFromFile(gssFileName)
	if err != nil {
		return nil, err
	}

	info.synch = synch
	info.src = src
	info.dst = dst
	info.items = items

	return info, nil
}

func populateInfoFromPath(gssFileName string, info *ScriptsInfo) error {
	base := filepath.Base(gssFileName)
	info.name = strings.TrimSuffix(base, path.Ext(base))

	dir := filepath.Dir(gssFileName)
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("unable to resolve absolute path for %s: %w", dir, err)
	}
	info.dir = absDir
	return nil
}

func readGSSConfigFromFile(gssFileName string) (string, string, string, []string, error) {
	gssFile, err := os.Open(gssFileName)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("unable to open %s: %w", gssFileName, err)
	}
	defer gssFile.Close()

	return readGSSConfig(gssFile, gssFileName)
}

func readGSSConfig(r io.Reader, label string) (string, string, string, []string, error) {
	input := bufio.NewScanner(r)

	synch, _ := scanLine(input)
	src, _ := scanLine(input)
	dst, _ := scanLine(input)

	// Skip blank line
	scanLine(input)

	dirs := mapset.NewSet[string]()
	for input.Scan() {
		d := strings.Trim(input.Text(), " /")
		if len(d) > 0 {
			dirs.Add(d)
		}
	}
	if err := input.Err(); err != nil {
		return "", "", "", nil, fmt.Errorf("unable to read %s: %w", label, err)
	}

	dirsSlice := dirs.ToSlice()
	sort.Strings(dirsSlice)

	return synch, src, dst, dirsSlice, nil
}

func scanLine(scanner *bufio.Scanner) (string, bool) {
	if !scanner.Scan() {
		return "", false
	}
	return scanner.Text(), true
}

// writeScripts orchestrates script generation for a parsed ScriptsInfo.
func (g *RsyncScriptGenerator) writeScripts(info *ScriptsInfo) error {
	g.printPlan(info)

	ext := ".sh"
	if g.Powershell {
		ext = ".ps1"
	}

	mainPath := filepath.Join(g.AutoGenDir, info.name+ext)
	if err := writeExecutableScript(mainPath, g.buildAllItemsScript(info)); err != nil {
		return fmt.Errorf("unable to write script %s: %w", mainPath, err)
	}

	if g.Files || len(info.items) <= 1 {
		return nil
	}

	return g.writePerItemScripts(info)
}

// writePerItemScripts creates one shell script per item in a subdirectory.
func (g *RsyncScriptGenerator) writePerItemScripts(info *ScriptsInfo) error {
	itemsDir := filepath.Join(g.AutoGenDir, info.name)
	if err := os.MkdirAll(itemsDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create directory %s: %w", itemsDir, err)
	}

	for _, item := range info.items {
		ext := ".sh"
		if g.Powershell {
			ext = ".ps1"
		}

		p := filepath.Join(itemsDir, item+ext)
		if err := writeExecutableScript(p, g.buildSingleItemScript(info, item)); err != nil {
			return fmt.Errorf("unable to write script %s: %w", p, err)
		}
	}
	return nil
}

// printPlan logs what's about to be generated.
func (g *RsyncScriptGenerator) printPlan(info *ScriptsInfo) {
	fmt.Printf("Generating scripts in %v\n", g.AutoGenDir)
	fmt.Printf("Synch root: %v\n", info.synch)
	fmt.Printf("Source: %v\n", info.src)
	fmt.Printf("Destination: %v\n", info.dst)

	label := "Directories"
	if g.Files {
		label = "Files"
	}
	fmt.Printf("%s to synch:\n", label)
	for _, item := range info.items {
		fmt.Println(item)
	}
	fmt.Println()
}

// --- Script building ---

func (g *RsyncScriptGenerator) buildAllItemsScript(info *ScriptsInfo) []byte {
	var b bytes.Buffer
	writeHeader(&b, g.Powershell)
	for _, item := range info.items {
		writeItemCommands(&b, g.Files, info, item, g.Powershell)
		b.WriteString("\n")
	}
	if g.Powershell {
		b.WriteString("\nGet-Date\n")
	} else {
		b.WriteString("\ndate\n")
	}
	return b.Bytes()
}

func (g *RsyncScriptGenerator) buildSingleItemScript(info *ScriptsInfo, item string) []byte {
	var b bytes.Buffer
	writeHeader(&b, g.Powershell)
	writeItemCommands(&b, false, info, item, g.Powershell)
	b.WriteString("\n")
	if g.Powershell {
		b.WriteString("\nGet-Date\n")
	} else {
		b.WriteString("\ndate\n")
	}
	return b.Bytes()
}

func writeHeader(b *bytes.Buffer, powershell bool) {
	if powershell {
		b.WriteString("#!/usr/bin/env pwsh\n# AUTOGEN'D - DO NOT EDIT!\n")
		fmt.Fprintf(b, "# Generated on %s\n\n", getNowFmt())
		b.WriteString("Get-Date\n\n")
	} else {
		b.WriteString("#!/bin/bash\n# AUTOGEN'D - DO NOT EDIT!\n")
		fmt.Fprintf(b, "# Generated on %s\n\n", getNowFmt())
		b.WriteString("date\n\n")
	}
}

func writeItemCommands(b *bytes.Buffer, files bool, info *ScriptsInfo, item string, powershell bool) {
	to := getCmdLine(files, info.synch, item, info.src, info.dst)
	fmt.Fprintf(b, "%s\n%s\n", getEchoLine(to, powershell), to)

	from := getCmdLine(files, info.synch, item, info.dst, info.src)
	fmt.Fprintf(b, "%s\n%s\n", getEchoLine(from, powershell), from)
}

// --- File helpers ---

func writeExecutableScript(scriptPath string, content []byte) error {
	if err := deleteIfExists(scriptPath); err != nil {
		return err
	}
	return os.WriteFile(scriptPath, content, 0o755)
}

func deleteIfExists(filePath string) error {
	_, err := os.Stat(filePath)
	if err == nil {
		fmt.Printf("%v exists - Deleting...\n", filePath)
		return os.Remove(filePath)
	}
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// --- Formatting helpers ---

func getCmdLine(files bool, synchRoot, item, src, dst string) string {
	if files {
		return fmt.Sprintf(filesCmdLineTemplate, synchRoot, src, item, dst)
	}
	return fmt.Sprintf(dirsCmdLineTemplate, synchRoot, src, item, dst, item)
}

func getEchoLine(cmd string, powershell bool) string {
	if powershell {
		return fmt.Sprintf("Write-Host '%s'", cmd)
	}
	return fmt.Sprintf("echo '%s'", cmd)
}

func getNowFmt() string {
	return time.Now().UTC().Format(time.RFC1123)
}
