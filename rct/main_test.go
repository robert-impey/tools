package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func TestReadTestFileErrorExt(t *testing.T) {
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "test.err")
	testContent := "command --fail\n---\nexpected error output"
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	command, output, err := readTestFile(testPath, filepath.Join(tempDir, "data"), false)
	if err != nil {
		t.Fatalf("readTestFile failed unexpectedly: %v", err)
	}

	if command != "command --fail" {
		t.Errorf("Expected command %q, got %q", "command --fail", command)
	}
	if output != "expected error output" {
		t.Errorf("Expected output %q, got %q", "expected error output", output)
	}
}

func TestReadTestFileReplacesTestDataDirToken(t *testing.T) {
	testDataDir := filepath.Join("/tmp", "fixture-data")
	testPath := filepath.Join(t.TempDir(), "test.txt")
	testContent := "echo TEST_DATA_DIR\n---\nexpected output"
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	command, output, err := readTestFile(testPath, testDataDir, false)
	if err != nil {
		t.Fatalf("readTestFile failed unexpectedly: %v", err)
	}

	if command != "echo "+testDataDir {
		t.Fatalf("command = %q, want %q", command, "echo "+testDataDir)
	}
	if output != "expected output" {
		t.Fatalf("output = %q, want %q", output, "expected output")
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

func TestRunTest(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	rootDir := t.TempDir()
	programDir := filepath.Join(rootDir, "program")
	testsDir := filepath.Join(programDir, testsDirName)
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		t.Fatal(err)
	}

	executableBytes, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(programDir, filepath.Base(executable)), executableBytes, 0755); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name           string
		extension      string
		mode           string
		outputType     string
		expectedOutput string
		wantSuccess    bool
	}{
		{name: "txt compares stdout", extension: ".txt", mode: "stdout", outputType: "out", expectedOutput: programDir, wantSuccess: true},
		{name: "err compares stderr", extension: ".err", mode: "stderr", outputType: "err", expectedOutput: programDir, wantSuccess: true},
		{name: "mismatched output fails", extension: ".txt", mode: "stdout", outputType: "out", expectedOutput: "unexpected output", wantSuccess: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("RCT_RUN_TEST_HELPER", tc.mode)

			testFile := filepath.Join(testsDir, "command"+tc.extension)
			content := fmt.Sprintf("%s -test.run=^TestRunTestCommandHelper$\n---\n%s",
				filepath.Base(executable), tc.expectedOutput)
			if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}

			got := runTest(context.Background(), testFile, rootDir, false, tc.outputType, programDir)
			if got != tc.wantSuccess {
				t.Errorf("runTest() = %t, want %t", got, tc.wantSuccess)
			}
		})
	}
}

func TestRunTestCommandHelper(t *testing.T) {
	switch os.Getenv("RCT_RUN_TEST_HELPER") {
	case "stdout":
		workingDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(workingDir)
		os.Exit(0)
	case "stderr":
		workingDir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintln(os.Stderr, workingDir)
		os.Exit(0)
	default:
		t.Skip("helper test is only run by TestRunTest")
	}
}

func TestMainUsesFlagsEndToEnd(t *testing.T) {
	if os.Getenv("RCT_CLI_HELPER_PROCESS") == "1" {
		os.Args = []string{"rct", os.Getenv("RCT_CLI_DIR"), "--test-data-directory", os.Getenv("RCT_CLI_TEST_DATA_DIR"), "--verbose"}
		main()
		return
	}

	rootDir := t.TempDir()
	programDir := filepath.Join(rootDir, "program")
	testsDir := filepath.Join(programDir, testsDirName)
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		t.Fatal(err)
	}

	testDataDir := filepath.Join(rootDir, "fixtures")
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatal(err)
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	binaryPath := filepath.Join(programDir, filepath.Base(executable))
	executableBytes, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binaryPath, executableBytes, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(testsDir, "cli.txt")
	content := fmt.Sprintf("%s -test.run=^TestMainHelperCommand$\n---\n%s", filepath.Base(executable), programDir)
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestMainUsesFlagsEndToEnd$")
	cmd.Env = append(os.Environ(),
		"RCT_CLI_HELPER_PROCESS=1",
		"RCT_CLI_DIR="+rootDir,
		"RCT_CLI_TEST_DATA_DIR="+testDataDir,
		"RCT_TEST_BINARY_HELPER=1",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI helper process exited unexpectedly: %v\nOutput:\n%s", err, output)
	}

	stdout := string(output)
	for _, needle := range []string{"Searching", "Program directory:", "Success rate:", "OK"} {
		if !strings.Contains(stdout, needle) {
			t.Fatalf("expected CLI output to contain %q, got:\n%s", needle, stdout)
		}
	}
	if !strings.Contains(stdout, programDir) {
		t.Fatalf("expected CLI output to include program dir %q, got:\n%s", programDir, stdout)
	}
}

func TestExecuteFindsAndRunsTestsAcrossProgramDirectories(t *testing.T) {
	rootDir := t.TempDir()
	programA := filepath.Join(rootDir, "programA")
	programB := filepath.Join(rootDir, "programB")
	for _, dir := range []string{programA, programB} {
		if err := os.MkdirAll(filepath.Join(dir, testsDirName), 0755); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.MkdirAll(filepath.Join(rootDir, ".ignored"), 0755); err != nil {
		t.Fatal(err)
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{programA, programB} {
		binaryPath := filepath.Join(dir, filepath.Base(executable))
		executableBytes, err := os.ReadFile(executable)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(binaryPath, executableBytes, 0755); err != nil {
			t.Fatal(err)
		}
	}

	for _, tc := range []struct {
		dir      string
		fileName string
		content  string
		expected string
		testType string
		exitCode int
		wantPass bool
	}{
		{dir: programA, fileName: "first.txt", content: fmt.Sprintf("%s -test.run=^TestExecuteHelperSuccess$\n---\n%s", filepath.Base(executable), programA), expected: programA, testType: "out", wantPass: true},
		{dir: programA, fileName: "second.err", content: fmt.Sprintf("%s -test.run=^TestExecuteHelperSuccess$\n---\n%s", filepath.Base(executable), programA), expected: programA, testType: "err", wantPass: true},
		{dir: programB, fileName: "third.txt", content: fmt.Sprintf("%s -test.run=^TestExecuteHelperSuccess$\n---\n%s", filepath.Base(executable), "wrong value"), expected: "wrong value", testType: "out", wantPass: false},
	} {
		if err := os.WriteFile(filepath.Join(tc.dir, testsDirName, tc.fileName), []byte(tc.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	settings := CommandSettings{Directory: rootDir, TestDataDirectory: rootDir, Verbose: false}
	code := execute(context.Background(), settings)
	if code != 0 {
		t.Fatalf("execute() returned %d, want 0", code)
	}

	if got := enumerateFilteredDirectories(rootDir); len(got) < 3 {
		t.Fatalf("enumerateFilteredDirectories() unexpectedly short: %v", got)
	}
}

func TestExecuteReturnsExpectedExitCodes(t *testing.T) {
	rootDir := t.TempDir()
	settings := CommandSettings{Directory: rootDir, TestDataDirectory: rootDir, Verbose: false}

	if code := execute(context.Background(), settings); code != 1 {
		t.Fatalf("execute() with no tests returned %d, want 1", code)
	}

	programDir := filepath.Join(rootDir, "program")
	if err := os.MkdirAll(filepath.Join(programDir, testsDirName), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(programDir, testsDirName, "empty.txt"), []byte("echo hi\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if code := execute(context.Background(), settings); code != 0 {
		t.Fatalf("execute() with a passing test returned %d, want 0", code)
	}
}

func TestExecuteVerboseOutputAndCancellation(t *testing.T) {
	rootDir := t.TempDir()
	programDir := filepath.Join(rootDir, "program")
	if err := os.MkdirAll(filepath.Join(programDir, testsDirName), 0755); err != nil {
		t.Fatal(err)
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	binaryPath := filepath.Join(programDir, filepath.Base(executable))
	executableBytes, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binaryPath, executableBytes, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(programDir, testsDirName, "verbose.txt")
	content := fmt.Sprintf("%s -test.run=^TestExecuteHelperSuccess$\n---\n%s", filepath.Base(executable), programDir)
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	settings := CommandSettings{Directory: rootDir, TestDataDirectory: rootDir, Verbose: true}
	if code := execute(ctx, settings); code != 1 {
		t.Fatalf("execute() after cancellation returned %d, want 1", code)
	}

	output := captureStdout(func() {
		settings = CommandSettings{Directory: rootDir, TestDataDirectory: rootDir, Verbose: true}
		code := execute(context.Background(), settings)
		if code != 0 {
			t.Fatalf("execute() in verbose mode returned %d, want 0", code)
		}
	})

	for _, needle := range []string{"Searching", "Program directory:", "Test file:", "Command output (Actual):", "Success rate:"} {
		if !strings.Contains(output, needle) {
			t.Fatalf("expected verbose output to contain %q, got:\n%s", needle, output)
		}
	}
}

func TestMainExecutesTxtAndErrFilesInNestedProgramFolders(t *testing.T) {
	for _, tc := range []struct {
		name       string
		fileExt    string
		helperMode string
	}{
		{name: "txt file", fileExt: ".txt", helperMode: "stdout"},
		{name: "err file", fileExt: ".err", helperMode: "stderr"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if os.Getenv("RCT_REAL_CLI_HELPER_PROCESS") == "1" {
				os.Args = []string{"rct", os.Getenv("RCT_REAL_CLI_DIR"), "--test-data-directory", os.Getenv("RCT_REAL_CLI_TEST_DATA_DIR")}
				main()
				return
			}

			rootDir := t.TempDir()
			programDir := filepath.Join(rootDir, "nested", "program")
			testsDir := filepath.Join(programDir, testsDirName)
			if err := os.MkdirAll(testsDir, 0755); err != nil {
				t.Fatal(err)
			}

			executable, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			binaryPath := filepath.Join(programDir, filepath.Base(executable))
			binaryBytes, err := os.ReadFile(executable)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(binaryPath, binaryBytes, 0755); err != nil {
				t.Fatal(err)
			}

			fileName := "nested-cli" + tc.fileExt
			content := fmt.Sprintf("%s -test.run=^TestMainHelperRealCliCommand$\n---\n%s", filepath.Base(executable), programDir)
			if err := os.WriteFile(filepath.Join(testsDir, fileName), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command(os.Args[0], "-test.run=^TestMainExecutesTxtAndErrFilesInNestedProgramFolders$")
			cmd.Env = append(os.Environ(),
				"RCT_REAL_CLI_HELPER_PROCESS=1",
				"RCT_REAL_CLI_DIR="+rootDir,
				"RCT_REAL_CLI_TEST_DATA_DIR="+rootDir,
				"RCT_REAL_CLI_HELPER_MODE="+tc.helperMode,
			)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("real CLI helper exited unexpectedly: %v\nOutput:\n%s", err, output)
			}

			stdout := string(output)
			if !strings.Contains(stdout, "Success rate:") {
				t.Fatalf("expected CLI output to report success rate, got:\n%s", stdout)
			}
			if !strings.Contains(stdout, "100.0%") {
				t.Fatalf("expected CLI output to report a successful run, got:\n%s", stdout)
			}
			if !strings.Contains(stdout, "Test file: nested-cli") {
				t.Fatalf("expected CLI output to show the nested test file, got:\n%s", stdout)
			}
		})
	}
}

func TestMainHelperRealCliCommand(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("RCT_REAL_CLI_HELPER_MODE") == "stdout" {
		fmt.Println(workingDir)
		os.Exit(0)
	}
	if os.Getenv("RCT_REAL_CLI_HELPER_MODE") == "stderr" {
		fmt.Fprintln(os.Stderr, workingDir)
		os.Exit(0)
	}
	t.Skip("helper only used by CLI execution tests")
}

func captureStdout(fn func()) string {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	fn()
	w.Close()
	output, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(output)
}

func TestExecuteHelperSuccess(t *testing.T) {
	if os.Getenv("RCT_EXECUTE_HELPER_SUCCESS") != "1" {
		t.Skip("helper test is only run by execute tests")
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(wd)
	os.Exit(0)
}

func TestMainHelperCommand(t *testing.T) {
	if os.Getenv("RCT_TEST_BINARY_HELPER") != "1" {
		t.Skip("helper test is only run by TestMainUsesFlagsEndToEnd")
	}

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(workingDir)
	os.Exit(0)
}

func TestEnumerateFilteredDirectoriesHandlesErrorsAndDSYM(t *testing.T) {
	tempDir := t.TempDir()

	validSubdir := filepath.Join(tempDir, "valid_subdir")
	if err := os.Mkdir(validSubdir, 0755); err != nil {
		t.Fatal(err)
	}

	dsymDir := filepath.Join(tempDir, "bundle.dSYM")
	if err := os.Mkdir(dsymDir, 0755); err != nil {
		t.Fatal(err)
	}

	unreadableDir := filepath.Join(tempDir, "unreadable_dir")
	if err := os.Mkdir(unreadableDir, 0000); err == nil {
		// Ensure permissions are restored so t.TempDir cleanup succeeds
		t.Cleanup(func() {
			os.Chmod(unreadableDir, 0755)
		})
	}

	results := enumerateFilteredDirectories(tempDir)

	for _, dir := range results {
		if strings.HasSuffix(dir, ".dSYM") {
			t.Errorf("expected .dSYM directory to be filtered out, got %s", dir)
		}
	}

	foundValid := false
	for _, dir := range results {
		if dir == validSubdir {
			foundValid = true
			break
		}
	}
	if !foundValid {
		t.Errorf("expected %s in results, but was missing", validSubdir)
	}
}

func TestReadTestFileErrorsAndVerbose(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Non-existent file
	nonExistent := filepath.Join(tempDir, "does_not_exist.txt")
	if _, _, err := readTestFile(nonExistent, tempDir, false); err == nil {
		t.Errorf("expected error reading non-existent file, got nil")
	}

	// 2. Malformed test files (< 3 lines)
	cases := []struct {
		name    string
		content string
	}{
		{name: "empty file", content: ""},
		{name: "one line", content: "echo hello"},
		{name: "two lines", content: "echo hello\n---"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, tc.name+".txt")
			if err := os.WriteFile(filePath, []byte(tc.content), 0644); err != nil {
				t.Fatal(err)
			}
			if _, _, err := readTestFile(filePath, tempDir, false); err == nil {
				t.Errorf("expected malformed test file error for %s, got nil", tc.name)
			}
		})
	}

	// 3. Verbose mode
	validFile := filepath.Join(tempDir, "valid.txt")
	if err := os.WriteFile(validFile, []byte("echo hello\n---\nhello\n"), 0644); err != nil {
		t.Fatal(err)
	}

	output := captureStdout(func() {
		cmd, exp, err := readTestFile(validFile, tempDir, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmd != "echo hello" || exp != "hello" {
			t.Fatalf("unexpected cmd or exp: %q, %q", cmd, exp)
		}
	})

	if !strings.Contains(output, "Command: echo hello") {
		t.Errorf("expected verbose output to contain 'Command: echo hello', got %q", output)
	}
}

func TestExecuteCommandEmptyAndNonExitErrors(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()

	// 1. Empty and whitespace-only command
	for _, cmd := range []string{"", "   ", "\t\n"} {
		out, code := executeCommand(ctx, cmd, tempDir, "out")
		if out != "" || code != -1 {
			t.Errorf("executeCommand(%q) = (%q, %d), want (%q, -1)", cmd, out, code, "")
		}
	}

	// 2. Non-exit error (command binary not found)
	out, code := executeCommand(ctx, "nonexistent_command_12345", tempDir, "out")
	if code != -1 {
		t.Errorf("executeCommand for nonexistent binary returned exit code %d, want -1 (output: %q)", code, out)
	}
}

func TestExecuteSkipsSubdirectoriesAndUnrecognizedFiles(t *testing.T) {
	t.Setenv("RCT_EXECUTE_HELPER_SUCCESS", "1")

	rootDir := t.TempDir()
	programDir := filepath.Join(rootDir, "program")
	testsDir := filepath.Join(programDir, testsDirName)
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		t.Fatal(err)
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	binaryPath := filepath.Join(programDir, filepath.Base(executable))
	executableBytes, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binaryPath, executableBytes, 0755); err != nil {
		t.Fatal(err)
	}

	// 1. Nested subdirectory inside tests/
	if err := os.Mkdir(filepath.Join(testsDir, "nested_subfolder"), 0755); err != nil {
		t.Fatal(err)
	}

	// 2. Unrecognized file extensions
	if err := os.WriteFile(filepath.Join(testsDir, "notes.md"), []byte("# Notes"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testsDir, "data.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}

	// 3. Valid test file
	testFile := filepath.Join(testsDir, "01_valid.txt")
	content := fmt.Sprintf("%s -test.run=^TestExecuteHelperSuccess$\n---\n%s", filepath.Base(executable), programDir)
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	output := captureStdout(func() {
		settings := CommandSettings{Directory: rootDir, TestDataDirectory: rootDir, Verbose: false}
		code := execute(context.Background(), settings)
		if code != 0 {
			t.Fatalf("execute() returned %d, want 0", code)
		}
	})

	if !strings.Contains(output, "Success rate: \033[32m1\033[0m/\033[1m1\033[0m") {
		t.Errorf("expected 1 test to be discovered and run, output: %s", output)
	}
}

func TestExecuteDeterministicAlphabeticalOrder(t *testing.T) {
	t.Setenv("RCT_EXECUTE_HELPER_SUCCESS", "1")

	rootDir := t.TempDir()
	programDir := filepath.Join(rootDir, "program")
	testsDir := filepath.Join(programDir, testsDirName)
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		t.Fatal(err)
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	binaryPath := filepath.Join(programDir, filepath.Base(executable))
	executableBytes, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binaryPath, executableBytes, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files in reverse or shuffled order
	for _, name := range []string{"03_test.txt", "01_test.txt", "02_test.txt"} {
		testFile := filepath.Join(testsDir, name)
		content := fmt.Sprintf("%s -test.run=^TestExecuteHelperSuccess$\n---\n%s", filepath.Base(executable), programDir)
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	output := captureStdout(func() {
		settings := CommandSettings{Directory: rootDir, TestDataDirectory: rootDir, Verbose: true}
		code := execute(context.Background(), settings)
		if code != 0 {
			t.Fatalf("execute() returned %d, want 0", code)
		}
	})

	pos1 := strings.Index(output, "Test file: 01_test.txt")
	pos2 := strings.Index(output, "Test file: 02_test.txt")
	pos3 := strings.Index(output, "Test file: 03_test.txt")

	if pos1 == -1 || pos2 == -1 || pos3 == -1 {
		t.Fatalf("all three test files should be present in output, got:\n%s", output)
	}
	if !(pos1 < pos2 && pos2 < pos3) {
		t.Errorf("expected tests to execute in alphabetical order: pos1=%d, pos2=%d, pos3=%d", pos1, pos2, pos3)
	}
}

func TestRunTestHandlesErrorsAndVerbose(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "missing.txt")

	// 1. Error reading test file in runTest
	if ok := runTest(context.Background(), nonExistentFile, tempDir, false, "out", tempDir); ok {
		t.Errorf("runTest() with non-existent file returned true, want false")
	}

	// 2. Verbose mode when test fails
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	programDir := filepath.Join(tempDir, "prog")
	if err := os.MkdirAll(programDir, 0755); err != nil {
		t.Fatal(err)
	}
	binaryPath := filepath.Join(programDir, filepath.Base(executable))
	executableBytes, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binaryPath, executableBytes, 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(programDir, "test.txt")
	content := fmt.Sprintf("%s -test.run=^TestExecuteHelperSuccess$\n---\nexpected_different_output", filepath.Base(executable))
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	output := captureStdout(func() {
		ok := runTest(context.Background(), testFile, tempDir, true, "out", programDir)
		if ok {
			t.Errorf("runTest() should have failed")
		}
	})

	for _, needle := range []string{"Test output (Expected):", "Command output (Actual):", "Exit code:", "FAIL"} {
		if !strings.Contains(output, needle) {
			t.Errorf("expected verbose runTest output to contain %q, got:\n%s", needle, output)
		}
	}
}
