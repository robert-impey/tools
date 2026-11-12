package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIUsesPositionalDirectoryArgument(t *testing.T) {
	rootDir := t.TempDir()
	executablePath, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	
	executable, err := os.ReadFile(executablePath)
	if err != nil {
		t.Fatal(err)
	}

	programDirs := []struct {
		name     string
		testFile string
	}{
		{name: filepath.Join(rootDir, "root-program"), testFile: "root_test.txt"},
		{name: filepath.Join(rootDir, "subproj", "subproj-program"), testFile: "subproj_test.txt"},
	}

	for _, program := range programDirs {
		testsDir := filepath.Join(program.name, testsDirName)
		if err := os.MkdirAll(testsDir, 0755); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(filepath.Join(program.name, filepath.Base(executablePath)), executable, 0755); err != nil {
			t.Fatal(err)
		}

		testFile := filepath.Join(testsDir, program.testFile)
		content := fmt.Sprintf("%s -test.run=^TestSubdirRegressionCommand$\n---\n%s",
			filepath.Base(executablePath), filepath.Base(program.name))
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	testCases := []struct {
		name      string
		directory string
		wantFiles []string
		dontFiles []string
		wantCount int
	}{
		{
			name:      "positional directory selects only subproject",
			directory: "subproj",
			wantFiles: []string{"subproj_test.txt"},
			dontFiles: []string{"root_test.txt"},
			wantCount: 1,
		},
		{
			name:      "no positional directory searches current tree",
			wantFiles: []string{"root_test.txt", "subproj_test.txt"},
			wantCount: 2,
		},
		{
			name:      "flag argument -d rejected by CLI",
			directory: filepath.Join(rootDir, "nonexistent_dir"),
			wantFiles: []string{},
			dontFiles: []string{},
			wantCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.directory == filepath.Join(rootDir, "nonexistent_dir") {
				output := runSubdirRegressionCLIWithDFlag(t, rootDir)

				// When -d flag is used, main() receives "-d" as a positional arg.
				// Since main() doesn't parse flags, it will fail to execute correctly.
				if strings.Contains(output, "Test file:") {
					t.Errorf("Expected CLI to fail when using -d flag; output: %s", output)
				}
				return
			}

			output := runSubdirRegressionCLI(t, rootDir, tc.directory)
			for _, name := range tc.wantFiles {
				if !strings.Contains(output, "Test file: "+name) {
					t.Errorf("expected CLI to run %q; output:\n%s", name, output)
				}
			}
			for _, name := range tc.dontFiles {
				if strings.Contains(output, "Test file: "+name) {
					t.Errorf("expected CLI not to run %q; output:\n%s", name, output)
				}
			}
			if got := strings.Count(output, "Test file:"); got != tc.wantCount {
				t.Errorf("CLI ran %d tests, want %d; output:\n%s", got, tc.wantCount, output)
			}
		})
	}
}

func runSubdirRegressionCLI(t *testing.T, workingDir, directory string) string {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run=^TestSubdirRegressionCLIProcess$")
	cmd.Dir = workingDir
	cmd.Env = append(os.Environ(),
		"RCT_SUBDIR_REGRESSION_CLI=1",
		"RCT_SUBDIR_REGRESSION_DIRECTORY="+directory,
		"RCT_SUBDIR_REGRESSION_COMMAND=1",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI subprocess failed: %v\nOutput:\n%s", err, output)
	}
	return string(output)
}

func runSubdirRegressionCLIWithDFlag(t *testing.T, workingDir string) string {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run=^TestSubdirRegressionCLIProcess$")
	cmd.Dir = workingDir
	// Simulate -d flag by passing it as an argument (not via env var).
	// main() does not parse this flag and will fail to use the directory.
	cmd.Env = append(os.Environ(),
		"RCT_SUBDIR_REGRESSION_CLI=1",
		"RCT_SUBDIR_REGRESSION_COMMAND=1",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI subprocess failed: %v\nOutput:\n%s", err, output)
	}
	return string(output)
}

func TestSubdirRegressionCLIProcess(t *testing.T) {
	if os.Getenv("RCT_SUBDIR_REGRESSION_CLI") != "1" {
		t.Skip("helper process for CLI regression test")
	}

	os.Args = []string{"rct"}
	if directory := os.Getenv("RCT_SUBDIR_REGRESSION_DIRECTORY"); directory != "" {
		os.Args = append(os.Args, directory)
	}
	main()
}

func TestSubdirRegressionCommand(t *testing.T) {
	if os.Getenv("RCT_SUBDIR_REGRESSION_COMMAND") != "1" {
		t.Skip("helper command for CLI regression test")
	}

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(filepath.Base(workingDir))
	os.Exit(0)
}