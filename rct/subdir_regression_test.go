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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

	for _, flag := range []string{"-d", "--directory"} {
		t.Run("rejects "+flag, func(t *testing.T) {
			output, err := runSubdirRegressionCLIWithFlag(t, rootDir, flag)
			if err == nil {
				t.Fatalf("Expected CLI to reject %q; output: %s", flag, output)
			}
			if strings.Contains(output, "Test file:") {
				t.Errorf("Expected CLI not to run tests when %q is used; output: %s", flag, output)
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

func runSubdirRegressionCLIWithFlag(t *testing.T, workingDir, flag string) (string, error) {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run=^TestSubdirRegressionCLIProcess$")
	cmd.Dir = workingDir
	cmd.Env = append(os.Environ(),
		"RCT_SUBDIR_REGRESSION_CLI=1",
		"RCT_SUBDIR_REGRESSION_FLAG="+flag,
		"RCT_SUBDIR_REGRESSION_DIRECTORY=subproj",
		"RCT_SUBDIR_REGRESSION_COMMAND=1",
	)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestSubdirRegressionCLIProcess(t *testing.T) {
	if os.Getenv("RCT_SUBDIR_REGRESSION_CLI") != "1" {
		t.Skip("helper process for CLI regression test")
	}

	os.Args = []string{"rct"}
	if flag := os.Getenv("RCT_SUBDIR_REGRESSION_FLAG"); flag != "" {
		os.Args = append(os.Args, flag)
	}
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