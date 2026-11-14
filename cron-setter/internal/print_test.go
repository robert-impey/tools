package cronsetter

import (
	"bytes"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestPrintAllTasksGeneratesValidCronEntries(t *testing.T) {
	t.Setenv("LOCAL_SCRIPTS", filepath.Join(t.TempDir(), "local-scripts"))

	var out bytes.Buffer
	PrintTasks(&out)

	expectedTools := []string{
		"Sweep-ManagedFolders.ps1",
		"run-nightly.sh",
		"reset-perms-linux.sh",
		"Build.ps1",
		"Clear-OldLogs.ps1",
		"Write-ManagedFoldersList.ps1",
	}
	foundTools := make(map[string]bool, len(expectedTools))
	entries := 0

	for _, line := range strings.Split(out.String(), "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 7 {
			t.Errorf("expected a standard five-field cron entry and executable/script paths, got %q", line)
			continue
		}

		minute, err := strconv.Atoi(fields[0])
		if err != nil || minute < 0 || minute > 59 {
			t.Errorf("minute field must be between 0 and 59, got %q", fields[0])
		}
		hour, err := strconv.Atoi(fields[1])
		if err != nil || hour < 0 || hour > 23 {
			t.Errorf("hour field must be between 0 and 23, got %q", fields[1])
		}
		if fields[2] != "*" || fields[3] != "*" || fields[4] != "*" {
			t.Errorf("expected standard five-field cron syntax, got %q", line)
		}
		if fields[5] != ZshPath && fields[5] != PwshPath {
			t.Errorf("expected an absolute executable path, got %q", fields[5])
		}
		if !filepath.IsAbs(fields[5]) || !filepath.IsAbs(fields[6]) {
			t.Errorf("expected absolute executable and script paths, got %q", line)
		}

		for _, tool := range expectedTools {
			if strings.HasSuffix(fields[6], tool) {
				foundTools[tool] = true
			}
		}
		entries++
	}

	if entries == 0 {
		t.Fatal("expected at least one cron entry")
	}
	for _, tool := range expectedTools {
		if !foundTools[tool] {
			t.Errorf("expected generated tasks to include %s", tool)
		}
	}
}

func TestPrintAllTasksUsesLocalScriptsEnvironmentOverride(t *testing.T) {
	localScripts := filepath.Join(t.TempDir(), "custom-scripts")
	t.Setenv("LOCAL_SCRIPTS", localScripts)

	var out bytes.Buffer
	PrintTasks(&out)

	for _, script := range []string{
		filepath.Join(localScripts, "_Common", "stay_deleted", "Sweep-ManagedFolders.ps1"),
		filepath.Join(localScripts, "_Common", "synch", "run-nightly.sh"),
		filepath.Join(localScripts, "_Common", "reset_perms", "reset-perms-linux.sh"),
		filepath.Join(localScripts, "_Common", "build", "Build.ps1"),
		filepath.Join(localScripts, "_Common", "logs_deleter", "Clear-OldLogs.ps1"),
		filepath.Join(localScripts, "_Common", "Write-ManagedFoldersList.ps1"),
	} {
		if !strings.Contains(out.String(), script) {
			t.Errorf("expected generated crontab to contain override path %q", script)
		}
	}
}

func TestGetLocalScriptsUsesHomeDirectoryWhenEnvironmentIsUnset(t *testing.T) {
	t.Setenv("LOCAL_SCRIPTS", "")

	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("user.Current() returned error: %v", err)
	}

	want := filepath.Join(currentUser.HomeDir, "local-scripts")
	if got := getLocalScripts(); got != want {
		t.Errorf("getLocalScripts() = %q, want %q", got, want)
	}
}
