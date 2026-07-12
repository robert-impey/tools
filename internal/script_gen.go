package internal

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

func WriteHeader(w io.Writer, headerLine string) error {
	if _, err := fmt.Fprintln(w, headerLine); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "# Written %s\n\n", time.Now().Format(time.RFC1123Z)); err != nil {
		return err
	}
	return nil
}

func WriteShebang(w io.Writer) error {
	if _, err := fmt.Fprintln(w, "#!/usr/bin/env pwsh"); err != nil {
		return err
	}
	return nil
}

func CleanFolderPathForLogName(path string) string {
	p := strings.TrimSpace(path)
	// normalize slashes
	p = strings.ReplaceAll(p, "/", "\\")
	// replace backslashes with underscores
	p = strings.ReplaceAll(p, "\\", "_")
	// replace @ (remote host indicator in rsync paths) with __
	p = strings.ReplaceAll(p, "@", "__")

	// remove invalid filename characters and tildes
	invalid := regexp.MustCompile(`[<>:\"/\\|?*~\x00-\x1F]`)
	p = invalid.ReplaceAllString(p, "")

	return strings.TrimLeft(p, "_")
}
