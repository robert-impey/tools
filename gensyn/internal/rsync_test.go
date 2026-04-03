package internal

/*
Copyright © 2025 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"bytes"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	common "github.com/robert-impey/tools/internal"
	"github.com/stretchr/testify/assert"
)

func TestParseGSSFile(t *testing.T) {
	scriptsInfo, err := ParseGSSFile("cleopatra.txt")

	assert.Nil(t, err)
	assert.NotNil(t, scriptsInfo)

	assert.Equal(t, scriptsInfo.name, "cleopatra")
	assert.Equal(t, scriptsInfo.synch,
		"rsync --update --recursive --verbose --times --iconv=utf8 --perms --exclude-from=/home/robert/local-scripts/_Common/synch/rsync-excluded.txt")
	assert.Equal(t, scriptsInfo.src, "robert@cleopatra.robertimpey.com:~")
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

	err := GenerateSynchScripts(false, outputDir, "cleopatra.txt")
	assert.Nil(t, err)

	scriptsFile := path.Join(outputDir, "cleopatra.ps1")
	if _, err := os.Stat(scriptsFile); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("scripts file %s doesn't exist", scriptsFile)
	}

	scriptsDir := path.Join(outputDir, "cleopatra")
	if _, err := os.Stat(scriptsDir); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("scripts dir %s doesn't exist", scriptsDir)
	}

	configScript := path.Join(scriptsDir, "config.ps1")
	if _, err := os.Stat(configScript); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("config script file %s doesn't exist", configScript)
	}
}

func TestGenerateFilesSynchScripts(t *testing.T) {
	outputDir := t.TempDir()

	err := GenerateSynchScripts(true, outputDir, "ssh-config.txt")
	assert.Nil(t, err)

	scriptsFile := path.Join(outputDir, "ssh-config.ps1")
	if _, err := os.Stat(scriptsFile); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("scripts file %s doesn't exist", scriptsFile)
	}
}

func TestWriteScriptsCharacterization(t *testing.T) {
	info := &ScriptsInfo{
		name:  "cleopatra",
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
		mainPath := filepath.Join(outputDir, "cleopatra.ps1")
		content, _ := os.ReadFile(mainPath)
		script := string(content)

		// Verifying the 'To' and 'From' lines for the first item 'config'
		// This matches your dirsCmdLineTemplate logic
		assert.Contains(t, script, "rsync --flags user@host:~/config/ /local/path/config")
		assert.Contains(t, script, "rsync --flags /local/path/config/ user@host:~/config")

		// 2. Verify Sub-script creation
		// The code creates a folder named after the script for individual items
		subPath := filepath.Join(outputDir, "cleopatra", "config.ps1")
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

		content, _ := os.ReadFile(filepath.Join(outputDir, "dots.ps1"))
		script := string(content)

		// In file mode, the item is only appended to the source
		assert.Contains(t, script, "rsync --files-flags /src/path/.bashrc /dst/path")
	})
}

func TestWriteHeaderShebang(t *testing.T) {
	var b bytes.Buffer
	_ = common.WriteShebang(&b)
	firstLine := strings.SplitN(b.String(), "\n", 2)[0]
	assert.Equal(t, "#!/usr/bin/env pwsh", firstLine)
}

func TestReadGSSConfigFromReader(t *testing.T) {
	gss := bytes.NewBufferString(`rsync --flags
user@host:~
/local/path

config
data
config
`)

	synch, src, dst, items, err := readGSSConfig(gss, "buffer")
	assert.Nil(t, err)

	assert.Equal(t, "rsync --flags", synch)
	assert.Equal(t, "user@host:~", src)
	assert.Equal(t, "/local/path", dst)
	assert.Equal(t, []string{"config", "data"}, items)
}

func TestReadGSSConfigFromReaderError(t *testing.T) {
	gss := &errorReader{err: errors.New("boom")}

	_, _, _, _, err := readGSSConfig(gss, "buffer")
	assert.NotNil(t, err)
}

type errorReader struct {
	err error
}

func (r *errorReader) Read(_ []byte) (int, error) {
	return 0, r.err
}
