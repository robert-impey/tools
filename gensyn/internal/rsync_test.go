package internal

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"errors"
	"os"
	"path"
	"path/filepath"
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

func TestWriteScriptsCharacterization(t *testing.T) {
	info := &ScriptsInfo{
		name:  "merneith",
		synch: "rsync --flags",
		src:   "user@host:~",
		dst:   "/local/path",
		items: []string{"config", "data"},
	}

	t.Run("Directory Mode Full Output", func(t *testing.T) {
		outputDir := t.TempDir()
		g := &RsyncScriptGenerator{Files: false, AutoGenDir: outputDir}
		err := g.writeScripts(info)
		assert.Nil(t, err)

		// 1. Verify Main Script Content
		mainPath := filepath.Join(outputDir, "merneith.sh")
		content, _ := os.ReadFile(mainPath)
		script := string(content)

		// Verifying the 'To' and 'From' lines for the first item 'config'
		// This matches your dirsCmdLineTemplate logic
		assert.Contains(t, script, "rsync --flags user@host:~/config/ /local/path/config")
		assert.Contains(t, script, "rsync --flags /local/path/config/ user@host:~/config")

		// 2. Verify Sub-script creation
		// The code creates a folder named after the script for individual items
		subPath := filepath.Join(outputDir, "merneith", "config.sh")
		subContent, err := os.ReadFile(subPath)
		assert.Nil(t, err, "Sub-script should exist for 'config'")
		assert.Contains(t, string(subContent), "rsync --flags user@host:~/config/ /local/path/config")
	})

	t.Run("File Mode Full Output", func(t *testing.T) {
		outputDir := t.TempDir()
		fileInfo := &ScriptsInfo{
			name:  "dots",
			synch: "rsync --files-flags",
			src:   "/src/path",
			dst:   "/dst/path",
			items: []string{".bashrc"},
		}

		g := &RsyncScriptGenerator{Files: true, AutoGenDir: outputDir}
		err := g.writeScripts(fileInfo)
		assert.Nil(t, err)

		content, _ := os.ReadFile(filepath.Join(outputDir, "dots.sh"))
		script := string(content)

		// In file mode, the item is only appended to the source
		assert.Contains(t, script, "rsync --files-flags /src/path/.bashrc /dst/path")
	})
}
