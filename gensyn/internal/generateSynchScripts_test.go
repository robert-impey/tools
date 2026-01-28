package internal

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGSSFile(t *testing.T) {
	scriptsInfo, err := ParseGSSFile("merneith.txt")

	assert.Nil(t, err)
	assert.NotNil(t, scriptsInfo)

	assert.Equal(t, scriptsInfo.name, "merneith")
	assert.Equal(t, scriptsInfo.synch,
		"rsync --update --recursive --verbose --times --iconv=utf8 --perms --exclude-from=/home/robert/local-scripts/_Common/synch/rsync-excluded.txt")
	assert.Equal(t, scriptsInfo.src, "robert@merneith.robertimpey.com:~")
	assert.Equal(t, scriptsInfo.dst, "/home/robert")
	assert.Equal(t, len(scriptsInfo.items), 4)
}

func TestParseGSSFileBadFile(t *testing.T) {
	scriptsInfo, err := ParseGSSFile("does-not-exist.txt")

	assert.NotNil(t, err)
	assert.Nil(t, scriptsInfo)
}

func TestGenerateDirectorySynchScripts(t *testing.T) {
	outputDir := t.TempDir()

	err := GenerateSynchScripts(false, outputDir, "merneith.txt")
	assert.Nil(t, err)

	scriptsFile := path.Join(outputDir, "merneith.sh")
	if _, err := os.Stat(scriptsFile); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("scripts file %s doesn't exist", scriptsFile)
	}

	scriptsDir := path.Join(outputDir, "merneith")
	if _, err := os.Stat(scriptsDir); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("scripts dir %s doesn't exist", scriptsDir)
	}

	configScript := path.Join(scriptsDir, "config.sh")
	if _, err := os.Stat(configScript); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("config script file %s doesn't exist", configScript)
	}
}

func TestGenerateFilesSynchScripts(t *testing.T) {
	outputDir := t.TempDir()

	err := GenerateSynchScripts(true, outputDir, "ssh-config.txt")
	assert.Nil(t, err)

	scriptsFile := path.Join(outputDir, "ssh-config.sh")
	if _, err := os.Stat(scriptsFile); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("scripts file %s doesn't exist", scriptsFile)
	}
}
