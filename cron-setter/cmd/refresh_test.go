package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRefreshCrontabPreservesBeforeAndAfter(t *testing.T) {
	input := `# Custom header task
5 0 * * * /usr/bin/echo before

# Start of crontab tasks created by cron-setter at 2023-01-01 00:00:00

10 0 * * * /old/task

# End of cron tasks created by cron-setter at 2023-01-01 00:00:00

# Custom footer task
55 23 * * * /usr/bin/echo after
`

	var out bytes.Buffer
	if err := refreshCrontab(strings.NewReader(input), &out); err != nil {
		t.Fatalf("refreshCrontab returned error: %v", err)
	}

	result := out.String()

	if strings.Contains(result, "/old/task") {
		t.Fatalf("expected old managed task to be replaced, got:\n%s", result)
	}
	if strings.Contains(result, "2023-01-01 00:00:00") {
		t.Fatalf("expected old timestamp to be replaced, got:\n%s", result)
	}

	expectedPrefix := `# Custom header task
5 0 * * * /usr/bin/echo before

# Start of crontab tasks created by cron-setter at `
	if !strings.HasPrefix(result, expectedPrefix) {
		t.Fatalf("expected output to start with %q, got:\n%s", expectedPrefix, result)
	}

	expectedSuffix := "\n\n# Custom footer task\n55 23 * * * /usr/bin/echo after\n"
	if !strings.HasSuffix(result, expectedSuffix) {
		t.Fatalf("expected output to end with %q, got:\n%s", expectedSuffix, result)
	}

	if !strings.Contains(result, footerCommentPrefix+" at ") {
		t.Fatalf("expected output to contain footer comment, got:\n%s", result)
	}
}

func TestRefreshCrontabRepeatedRefreshDoesNotAccumulateBlankLines(t *testing.T) {
	input := `5 0 * * * /usr/bin/echo before
# Start of crontab tasks created by cron-setter at 2023-01-01 00:00:00
10 0 * * * /old/task
# End of cron tasks created by cron-setter at 2023-01-01 00:00:00
55 23 * * * /usr/bin/echo after
`

	var firstOut bytes.Buffer
	if err := refreshCrontab(strings.NewReader(input), &firstOut); err != nil {
		t.Fatalf("first refreshCrontab error: %v", err)
	}

	var secondOut bytes.Buffer
	if err := refreshCrontab(strings.NewReader(firstOut.String()), &secondOut); err != nil {
		t.Fatalf("second refreshCrontab error: %v", err)
	}

	firstLines := strings.Split(firstOut.String(), "\n")
	secondLines := strings.Split(secondOut.String(), "\n")

	if len(firstLines) != len(secondLines) {
		t.Fatalf("expected line count to remain %d on repeated refresh, got %d", len(firstLines), len(secondLines))
	}

	if firstLines[0] != "5 0 * * * /usr/bin/echo before" || secondLines[0] != "5 0 * * * /usr/bin/echo before" {
		t.Fatalf("unexpected first line after repeated refresh: %q", secondLines[0])
	}

	if firstLines[len(firstLines)-2] != "55 23 * * * /usr/bin/echo after" ||
		secondLines[len(secondLines)-2] != "55 23 * * * /usr/bin/echo after" {
		t.Fatalf("unexpected last non-empty line after repeated refresh: %q", secondLines[len(secondLines)-2])
	}
}

func TestRefreshCrontabWithoutExistingSectionAppendsTasks(t *testing.T) {
	input := "# Only custom task\n5 0 * * * /usr/bin/echo custom\n"

	var out bytes.Buffer
	if err := refreshCrontab(strings.NewReader(input), &out); err != nil {
		t.Fatalf("refreshCrontab returned error: %v", err)
	}

	result := out.String()
	if !strings.HasPrefix(result, input+headerCommentPrefix+" at ") {
		t.Fatalf("expected custom tasks followed by cron-setter header, got:\n%s", result)
	}
	if !strings.Contains(result, footerCommentPrefix+" at ") {
		t.Fatalf("expected output to contain footer comment, got:\n%s", result)
	}
}
