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

	common "github.com/robert-impey/tools/internal"
)

const filesCmdLineTemplate = "%v %v/%v %v/"
const dirsCmdLineTemplate = "%v %v/%v/ %v/%v"

// ScriptsInfo holds parsed data from a .gss file.
type ScriptsInfo struct {
	name, dir, synch, src, dst string
	items                      []string
}

// RsyncScriptGenerator holds the configuration for generating rsync scripts.
type RsyncScriptGenerator struct {
	Files      bool   // true = file mode, false = directory mode
	AutoGenDir string // output directory for generated scripts
}

// GenerateSynchScripts is a convenience wrapper preserving the original API.
func GenerateSynchScripts(files bool, autoGenDir string, gssFile string) error {
	g := &RsyncScriptGenerator{Files: files, AutoGenDir: autoGenDir}
	return g.GenerateSynchScripts(gssFile)
}

// GenerateSynchScripts parses a .gss file and writes rsync PowerShell scripts.
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

	dirs := make(map[string]struct{})

	for input.Scan() {
		d := strings.Trim(input.Text(), " /")
		if len(d) > 0 {
			dirs[d] = struct{}{}
		}
	}

	if err := input.Err(); err != nil {
		return "", "", "", nil, fmt.Errorf("unable to read %s: %w", label, err)
	}

	// convert to slice
	dirsSlice := make([]string, 0, len(dirs))
	for d := range dirs {
		dirsSlice = append(dirsSlice, d)
	}

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

	ext := ".ps1"
	mainPath := filepath.Join(g.AutoGenDir, info.name+ext)
	if err := writeExecutableScript(mainPath, g.buildAllItemsScript(info)); err != nil {
		return fmt.Errorf("unable to write script %s: %w", mainPath, err)
	}

	if g.Files || len(info.items) <= 1 {
		return nil
	}

	return g.writePerItemScripts(info)
}

// writePerItemScripts creates one PowerShell script per item in a subdirectory.
func (g *RsyncScriptGenerator) writePerItemScripts(info *ScriptsInfo) error {
	itemsDir := filepath.Join(g.AutoGenDir, info.name)
	if err := os.MkdirAll(itemsDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create directory %s: %w", itemsDir, err)
	}

	for _, item := range info.items {
		ext := ".ps1"
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
	writeRsyncHeader(&b, info)
	for _, item := range info.items {
		writeItemCommands(&b, g.Files, info, item)
		b.WriteString("\n")
	}
	b.WriteString("Get-Date\n")
	return b.Bytes()
}

func (g *RsyncScriptGenerator) buildSingleItemScript(info *ScriptsInfo, item string) []byte {
	var b bytes.Buffer
	writeRsyncHeader(&b, info)
	writeItemCommands(&b, false, info, item)
	b.WriteString("\n")
	b.WriteString("Get-Date\n")
	return b.Bytes()
}

func writeRsyncHeader(b *bytes.Buffer, info *ScriptsInfo) {
	_ = common.WriteShebang(b)
	_ = common.WriteHeader(b, "# AUTOGEN'D - DO NOT EDIT!")

	b.WriteString("param(\n")
	b.WriteString("    [Parameter (Mandatory = $False)]\n")
	b.WriteString("    [switch]$logged = $False\n")
	b.WriteString(")\n\n")

	b.WriteString("Get-Date\n\n")

	fmt.Fprintf(b, "$id = \"%s\"\n", info.name)
	fmt.Fprintf(b, "$srcLogName = \"%s\"\n", cleanFolderPathForLogName(info.src))
	fmt.Fprintf(b, "$dstLogName = \"%s\"\n", cleanFolderPathForLogName(info.dst))
	b.WriteString("\n")

	b.WriteString(sshVariablesBlock)
	b.WriteString("\n")
}

const sshVariablesBlock = `try {
    $hostname = [System.Net.Dns]::GetHostByName((hostname)).HostName
}
catch {
    $hostname = [System.Environment]::MachineName
}

if ($hostname -match "^[^.]+") {
    $hostname = $Matches[0]
}

@(
    "$env:HOME/.keychain/$($hostname)-sh"
    "$env:HOME/.ssh/environment-$($hostname)"
) | ForEach-Object {
    $keychainFile = $_

    if (Test-Path $keychainFile) {
        Write-Output "Loading $keychainFile"
        foreach ($line in Get-Content $keychainFile) {
            if ($line -match "SSH_AUTH_SOCK=([^;]+);.*") {
                $env:SSH_AUTH_SOCK = $Matches[1]

                Write-Output "SSH_AUTH_SOCK: $env:SSH_AUTH_SOCK"
            }

            if ($line -match "SSH_AGENT_PID=(\d+);.*") {
                $env:SSH_AGENT_PID = $Matches[1]

                Write-Output "SSH_AGENT_PID: $env:SSH_AGENT_PID"
            }
        }
    }
    else {
        Write-Output "$keychainFile does not exist!"
    }
}
`

func writeItemCommands(b *bytes.Buffer, files bool, info *ScriptsInfo, item string) {
	fmt.Fprintf(b, "$item = \"%s\"\n\n", item)

	to := getCmdLine(files, info.synch, item, info.src, info.dst)
	writeDirectionalCommand(b, to, "src", "dst")

	from := getCmdLine(files, info.synch, item, info.dst, info.src)
	writeDirectionalCommand(b, from, "dst", "src")
}

func writeDirectionalCommand(b *bytes.Buffer, cmd, fromLogVar, toLogVar string) {
	fmt.Fprintf(b, "%s\n", getEchoLine(cmd))
	fmt.Fprintln(b, "if ($logged)")
	fmt.Fprintln(b, "{")
	fmt.Fprintln(b, "    $logTimeStr = Get-Date -Format \"yyyy-MM-ddTHH_mm_ss\"")
	fmt.Fprintf(b, "    $logFileBase = \"$($logTimeStr).$($id).$($item).$($%sLogName)-to-$($%sLogName)\"\n", fromLogVar, toLogVar)
	fmt.Fprintln(b)
	fmt.Fprintln(b, "    $homeDir = $null")
	fmt.Fprintln(b, "    if ($env:USERPROFILE -and (Test-Path $env:USERPROFILE)) {")
	fmt.Fprintln(b, "        $homeDir = $env:USERPROFILE")
	fmt.Fprintln(b, "    } elseif ($env:HOME -and (Test-Path $env:HOME)) {")
	fmt.Fprintln(b, "        $homeDir = $env:HOME")
	fmt.Fprintln(b, "    }")
	fmt.Fprintln(b, "    $logsDir = Join-Path $homeDir \"logs\"")
	fmt.Fprintln(b, "    if (-not (Test-Path $logsDir)) {")
	fmt.Fprintln(b, "        New-Item -ItemType Directory -Path $logsDir -Force | Out-Null")
	fmt.Fprintln(b, "    }")
	fmt.Fprintln(b, "    $synchLogsDir = Join-Path $logsDir \"synch\"")
	fmt.Fprintln(b, "    if (-not (Test-Path $synchLogsDir)) {")
	fmt.Fprintln(b, "        New-Item -ItemType Directory -Path $synchLogsDir -Force | Out-Null")
	fmt.Fprintln(b, "    }")
	fmt.Fprintln(b, "    $logPathBase = Join-Path $synchLogsDir \"$($logFileBase).rsync-synch\"")
	fmt.Fprintln(b, "    $logFile = \"$($logPathBase).log\"")
	fmt.Fprintln(b, "    $errFile = \"$($logPathBase).err\"")
	fmt.Fprintf(b, "    %s 1> $logFile 2> $errFile\n", cmd)
	fmt.Fprintln(b, "} else")
	fmt.Fprintln(b, "{")
	fmt.Fprintf(b, "    %s\n", cmd)
	fmt.Fprintln(b, "}")
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

func getEchoLine(cmd string) string {
	return fmt.Sprintf("Write-Host '%s'", cmd)
}
