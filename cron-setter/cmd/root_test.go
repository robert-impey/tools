package cmd

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func resetRootCmd(t *testing.T) {
	t.Helper()
	prevIn := rootCmd.InOrStdin()
	prevOut := rootCmd.OutOrStdout()
	prevErr := rootCmd.ErrOrStderr()

	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetIn(prevIn)
		rootCmd.SetOut(prevOut)
		rootCmd.SetErr(prevErr)
	})
}

func TestRootCmdPrint(t *testing.T) {
	resetRootCmd(t)
	t.Setenv("LOCAL_SCRIPTS", filepath.Join(t.TempDir(), "local-scripts"))

	var out bytes.Buffer
	rootCmd.SetArgs([]string{"print"})
	rootCmd.SetOut(&out)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute(print) failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, headerCommentPrefix+" at ") {
		t.Errorf("expected output to contain header comment prefix, got: %s", result)
	}
	if !strings.Contains(result, footerCommentPrefix+" at ") {
		t.Errorf("expected output to contain footer comment prefix, got: %s", result)
	}
	if !strings.Contains(result, "Sweep-ManagedFolders.ps1") {
		t.Errorf("expected output to contain Sweep-ManagedFolders.ps1, got: %s", result)
	}
}

func TestRootCmdRefresh(t *testing.T) {
	resetRootCmd(t)
	t.Setenv("LOCAL_SCRIPTS", filepath.Join(t.TempDir(), "local-scripts"))

	input := `# Custom header
5 0 * * * /usr/bin/echo before

# Start of crontab tasks created by cron-setter at 2023-01-01 00:00:00
10 0 * * * /old/task
# End of cron tasks created by cron-setter at 2023-01-01 00:00:00

# Custom footer
55 23 * * * /usr/bin/echo after
`

	var out bytes.Buffer
	rootCmd.SetArgs([]string{"refresh"})
	rootCmd.SetIn(strings.NewReader(input))
	rootCmd.SetOut(&out)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute(refresh) failed: %v", err)
	}

	result := out.String()
	if strings.Contains(result, "/old/task") {
		t.Errorf("expected /old/task to be removed from refreshed crontab, got: %s", result)
	}
	if strings.Contains(result, "2023-01-01 00:00:00") {
		t.Errorf("expected old timestamp to be updated, got: %s", result)
	}
	if !strings.HasPrefix(result, "# Custom header\n5 0 * * * /usr/bin/echo before\n\n# Start of crontab tasks created by cron-setter at ") {
		t.Errorf("expected custom header preserved before regenerated block, got: %s", result)
	}
	if !strings.HasSuffix(result, "\n\n# Custom footer\n55 23 * * * /usr/bin/echo after\n") {
		t.Errorf("expected custom footer preserved after regenerated block, got: %s", result)
	}
}

func TestRootCmdHelp(t *testing.T) {
	resetRootCmd(t)

	var out bytes.Buffer
	rootCmd.SetArgs([]string{"--help"})
	rootCmd.SetOut(&out)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd.Execute(--help) failed: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "Usage:") {
		t.Errorf("expected help output to contain Usage:, got: %s", result)
	}
	if !strings.Contains(result, "print") || !strings.Contains(result, "refresh") {
		t.Errorf("expected help output to list subcommands print and refresh, got: %s", result)
	}
}

func TestExecuteEntrypoint(t *testing.T) {
	resetRootCmd(t)

	rootCmd.SetArgs([]string{"--help"})
	rootCmd.SetOut(io.Discard)

	Execute()
}
