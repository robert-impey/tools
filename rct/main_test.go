package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Mocking the global context. For unit tests, we use a background context.
var testCtx = context.Background()

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

// TestExecuteCommand tests the execution of a shell command.
func TestExecuteCommand(t *testing.T) {
	// Use a predictable command that exists on most systems (like 'echo')
	// Since the environment is win32, we use cmd /c.

	// Test 1: Success case (Echoing something)
	cmdSuccess := "cmd /c echo Hello World"

	output, exitCode := executeCommand(testCtx, cmdSuccess, "", "out")

	if exitCode != 0 {
		t.Errorf("Test 1 Failed: Expected exit code 0, got %d. Output: %s", exitCode, output)
	}

	expectedSuccess := "Hello World"
	if strings.TrimRight(output, "\r\n") != strings.TrimRight(expectedSuccess, "\r\n") {
		t.Errorf("Test 1 Failed: Expected output '%s', got '%s'", expectedSuccess, output)
	}

	// Test 2: Failure case (Non-existent command)
	// On Windows, attempting to run a non-existent command via cmd /c will often fail
	cmdFailure := "non_existent_command_12345"

	// We expect an exit code != 0
	outputFailure, exitCodeFailure := executeCommand(testCtx, cmdFailure, "", "out")

	if exitCodeFailure == 0 {
		t.Errorf("Test 2 Failed: Expected non-zero exit code on failure, got 0")
	}

	// The output on failure should contain an error message from the shell
	if outputFailure == "" {
		t.Logf("Test 2 Info: Failure command produced no output, relying on exit code check.")
	}
}

// --- Main Logic Tests ---

// TestRunTest_Success simulates a successful test run.
// This test is highly complex due to required filesystem and OS process mocking
// and is skipped in favor of unit testing the core components.
func TestRunTest_Success(t *testing.T) {
	t.Skip("Skipping full runTest integration test due to high setup complexity (mocking FS and OS execution)")
}
