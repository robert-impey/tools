package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// CommandSettings represents the CLI arguments
type CommandSettings struct {
	Directory         string
	TestDataDirectory string
	Verbose           bool
}

const testsDirName = "tests"

func main() {
	// Simple argument handling for demonstration
	settings := CommandSettings{
		Directory:         ".",
		TestDataDirectory: "./data",
		Verbose:           true,
	}

	ctx := context.Background() // Use for cancellation if needed
	exitCode := execute(ctx, settings)
	os.Exit(exitCode)
}

func execute(ctx context.Context, settings CommandSettings) int {
	dir, _ := filepath.Abs(settings.Directory)
	testDataDir, _ := filepath.Abs(settings.TestDataDirectory)

	successes := 0
	testsCount := 0

	if settings.Verbose {
		fmt.Printf("Searching %s\n", dir)
	}

	programDirs := enumerateFilteredDirectories(dir)

	for _, programDir := range programDirs {
		// Check for cancellation
		select {
		case <-ctx.Done():
			fmt.Println("\n\033[33mOperation cancelled.\033[0m")
			return 1
		default:
		}

		testDirPath := filepath.Join(programDir, testsDirName)

		if settings.Verbose {
			printSeparator('+', 40)
			fmt.Printf("\033[33mProgram directory:\033[0m \033[1m%s\033[0m\n", programDir)
		}

		if info, err := os.Stat(testDirPath); err == nil && info.IsDir() {
			files, _ := os.ReadDir(testDirPath)

			// Sort files by name for consistency
			sort.Slice(files, func(i, j int) bool {
				return files[i].Name() < files[j].Name()
			})

			for _, file := range files {
				if file.IsDir() {
					continue
				}

				ext := strings.ToLower(filepath.Ext(file.Name()))
				var testType string
				switch ext {
				case ".txt":
					testType = "out"
				case ".err":
					testType = "err"
				default:
					continue
				}

				testsCount++
				testFile := filepath.Join(testDirPath, file.Name())
				if runTest(ctx, testFile, testDataDir, settings.Verbose, testType, programDir) {
					successes++
				}
			}
		}

		if settings.Verbose {
			printSeparator('+', 40)
		}
	}

	if testsCount == 0 {
		fmt.Println("\033[31mNo tests found!\033[0m")
		return 1
	}

	successRate := float64(successes) / float64(testsCount) * 100.0
	fmt.Printf("Success rate: \033[32m%d\033[0m/\033[1m%d\033[0m (\033[32m%.1f%%\033[0m)\n",
		successes, testsCount, successRate)

	return 0
}

func enumerateFilteredDirectories(rootDir string) []string {
	var programDirs []string
	queue := []string{rootDir}
	programDirs = append(programDirs, rootDir)

	for len(queue) > 0 {
		currentDir := queue[0]
		queue = queue[1:]

		entries, err := os.ReadDir(currentDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				fullPath := filepath.Join(currentDir, name)

				if !strings.HasPrefix(name, ".") && !strings.HasSuffix(name, ".dSYM") {
					programDirs = append(programDirs, fullPath)
					queue = append(queue, fullPath)
				}
			}
		}
	}
	return programDirs
}

func runTest(ctx context.Context, testFile, testDataDir string, verbose bool, testType, programDir string) bool {
	if verbose {
		printSeparator('-', 40)
	}

	fileName := filepath.Base(testFile)
	fmt.Printf("Test file: %s", fileName)
	if verbose {
		fmt.Println()
	} else {
		fmt.Print(" ")
	}

	command, expectedOutput, err := readTestFile(testFile, testDataDir, verbose)
	if err != nil {
		fmt.Printf("\033[31mError reading test: %v\033[0m\n", err)
		return false
	}

	if verbose {
		fmt.Println("Test output (Expected):")
		printSeparator('.', 40)
		fmt.Println(expectedOutput)
		printSeparator('.', 40)
	}

	actualOutput, exitCode := executeCommand(ctx, command, programDir, testType)

	if verbose {
		fmt.Println("Command output (Actual):")
		printSeparator('.', 40)
		fmt.Println(actualOutput)
		printSeparator('.', 40)
		fmt.Printf("Exit code: %d\n", exitCode)
		printSeparator('.', 40)
	}

	success := strings.TrimRight(actualOutput, "\r\n") == strings.TrimRight(expectedOutput, "\r\n")

	if success {
		fmt.Println("\033[32mOK\033[0m")
	} else {
		fmt.Println("\033[31mFAIL\033[0m")
	}

	return success
}

func readTestFile(path, testDataDir string, verbose bool) (string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) < 3 {
		return "", "", fmt.Errorf("malformed test file")
	}

	command := strings.TrimSpace(lines[0])
	command = strings.ReplaceAll(command, "TEST_DATA_DIR", testDataDir)

	if verbose {
		fmt.Printf("Command: %s\n", command)
	}

	expectedOutput := strings.Join(lines[2:], "\n")
	return command, expectedOutput, nil
}

func executeCommand(ctx context.Context, commandStr, workingDir, outputType string) (string, int) {
	parts := strings.Fields(commandStr)
	if len(parts) == 0 {
		return "", -1
	}

	executable := parts[0]
	args := parts[1:]

	// Use full path for executable
	cmdPath := filepath.Join(workingDir, executable)
	cmd := exec.CommandContext(ctx, cmdPath, args...)
	cmd.Dir = workingDir

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1
		}
	}

	var output string
	if outputType == "out" {
		output = stdout.String()
	} else {
		output = stderr.String()
	}

	return strings.TrimRight(output, "\r\n"), exitCode
}

func printSeparator(char rune, count int) {
	line := strings.Repeat(string(char), count)
	fmt.Printf("\n%s\n\n", line)
}
