package internal

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSearchDirectory_OutputRedirection(t *testing.T) {
	// Setup a temporary directory for searching
	searchDir := t.TempDir()

	// Create some files that will trigger matching stems
	// a.txt and a1.txt should match (assuming FindMatchingStems logic)
	// Actually looking at FindMatchingStems in tidy_folder.go:
	// It matches files with same extension where one stem is prefix of another.
	err := os.WriteFile(filepath.Join(searchDir, "a.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(searchDir, "a1.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Redirection to stdout when logsDir is empty", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := SearchDirectory(searchDir, "")

		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("SearchDirectory failed: %v", err)
		}

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		output := buf.String()

		if !strings.Contains(output, "--- Results for:") {
			t.Errorf("Expected output to contain '--- Results for:', but got: %q", output)
		}
	})

	t.Run("Redirection to file when logsDir is set", func(t *testing.T) {
		logsDir := t.TempDir()

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := SearchDirectory(searchDir, logsDir)

		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("SearchDirectory failed: %v", err)
		}

		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		stdoutOutput := buf.String()

		// The issue says it currently goes to BOTH stdout and the log file.
		// We want it to NOT go to stdout (except for the "OK: processed" message maybe,
		// but the core results "--- Results for:" should only be in the log).
		// Looking at the code:
		// if logFile != nil && funcErr == nil { fmt.Printf("OK: processed %s\n", dir) }
		// This "OK" message might be acceptable on stdout, but the "Results for" should not be.

		if strings.Contains(stdoutOutput, "--- Results for:") {
			t.Errorf("Expected '--- Results for:' NOT to be in stdout when logsDir is set, but it was found. Output: %q", stdoutOutput)
		}

		// Check if log file exists and contains the results
		files, err := os.ReadDir(logsDir)
		if err != nil {
			t.Fatal(err)
		}

		foundLog := false
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".log") {
				foundLog = true
				content, err := os.ReadFile(filepath.Join(logsDir, f.Name()))
				if err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(string(content), "--- Results for:") {
					t.Errorf("Log file %s does not contain '--- Results for:'. Content: %q", f.Name(), string(content))
				}
			}
		}

		if !foundLog {
			t.Error("No log file was created in logsDir")
		}
	})
}
