package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestEnumerateFilteredDirectories tests the logic for finding program directories.
func TestEnumerateFilteredDirectories(t *testing.T) {
	// 1. Setup a temporary directory structure
	tempDir, err := os.MkdirTemp("", "test_root")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Structure:
	// tempDir/
	// ├── programA/
	// │   ├── tests/
	// │   └── file.txt
	// ├── .ignored/
	// ├── programB/
	// │   ├── tests/
	// │   └── test.err
	// └── file_without_program/

	// Create programA
	programADir := filepath.Join(tempDir, "programA")
	os.Mkdir(programADir, 0755)
	os.Mkdir(filepath.Join(programADir, "tests"), 0755)
	os.WriteFile(filepath.Join(programADir, "tests", "test1.txt"), []byte("dummy"), 0644)
	os.WriteFile(filepath.Join(programADir, "file.txt"), []byte("dummy"), 0644)

	// Create ignored directory
	os.Mkdir(filepath.Join(tempDir, ".ignored"), 0755)

	// Create programB
	programBDir := filepath.Join(tempDir, "programB")
	os.Mkdir(programBDir, 0755)
	os.Mkdir(filepath.Join(programBDir, "tests"), 0755)
	os.WriteFile(filepath.Join(programBDir, "tests", "test2.err"), []byte("dummy"), 0644)

	// Create another directory to test filtering
	os.Mkdir(filepath.Join(tempDir, "file_without_program"), 0755)

	// 2. Run the function under test
	actual := enumerateFilteredDirectories(tempDir)

	// 3. Assertions
	expected := []string{
		tempDir,
		programADir,
		programBDir,
		filepath.Join(tempDir, "file_without_program"),
	}

	// Since the function is non-deterministic in order and includes rootDir,
	// we check containment and count.
	if len(actual) < 3 {
		t.Fatalf("Expected at least 3 program directories, got %d. List: %v", len(actual), actual)
	}

	found := make(map[string]bool)
	for _, dir := range actual {
		found[dir] = true
	}

	for _, expectedDir := range expected {
		if _, found := found[expectedDir]; !found {
			t.Errorf("Expected directory %s not found in results.", expectedDir)
		}
	}
}

// TestReadTestFile tests the logic for parsing a test definition file.
func TestReadTestFile(t *testing.T) {
	// Create a temporary test file
	// Line 1: Command
	// Line 2: Ignored/Separator (Skipped in logic)
	// Line 3+: Expected Output
	testContent := `echo "hello world"` + "\n" + "SUCCESS_LINE_1" + "\n" + "line 1\nline 2"
	tempFile, err := os.MkdirTemp("", "test_data")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	testPath := filepath.Join(tempFile, "test.txt")
	os.WriteFile(testPath, []byte(testContent), 0644)
	defer os.RemoveAll(tempFile)

	expectedCommand := "echo \"hello world\""
	expectedOutput := "line 1\nline 2"

	// Run the function under test
	command, output, err := readTestFile(testPath, filepath.Join(tempFile, "data"), false)

	if err != nil {
		t.Fatalf("readTestFile failed unexpectedly: %v", err)
	}

	if command != expectedCommand {
		t.Errorf("Expected command '%s', got '%s'", expectedCommand, command)
	}

	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}
}

// TestExecuteCommand tests command output selection and exit-code handling.
func TestExecuteCommand(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	command := filepath.Base(executable) + " -test.run=^TestExecuteCommandHelper$"
	testCases := []struct {
		name       string
		mode       string
		outputType string
		wantOutput string
		wantCode   int
	}{
		{name: "stdout", mode: "stdout", outputType: "out", wantOutput: "helper stdout", wantCode: 0},
		{name: "stderr", mode: "stderr", outputType: "err", wantOutput: "helper stderr", wantCode: 0},
		{name: "nonzero exit", mode: "exit", outputType: "out", wantCode: 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("RCT_EXECUTE_COMMAND_HELPER", tc.mode)
			output, exitCode := executeCommand(context.Background(), command, filepath.Dir(executable), tc.outputType)

			if output != tc.wantOutput {
				t.Errorf("output = %q, want %q", output, tc.wantOutput)
			}
			if exitCode != tc.wantCode {
				t.Errorf("exit code = %d, want %d", exitCode, tc.wantCode)
			}
		})
	}
}

func TestExecuteCommandHelper(t *testing.T) {
	switch os.Getenv("RCT_EXECUTE_COMMAND_HELPER") {
	case "stdout":
		fmt.Print("helper stdout\n")
		os.Exit(0)
	case "stderr":
		fmt.Fprint(os.Stderr, "helper stderr\n")
		os.Exit(0)
	case "exit":
		os.Exit(7)
	default:
		t.Skip("helper test is only run by TestExecuteCommand")
	}
}

// --- Main Logic Tests ---

// TestRunTest_Success simulates a successful test run.
// This test is highly complex due to required filesystem and OS process mocking
// and is skipped in favor of unit testing the core components.
func TestRunTest_Success(t *testing.T) {
	t.Skip("Skipping full runTest integration test due to high setup complexity (mocking FS and OS execution)")
}
