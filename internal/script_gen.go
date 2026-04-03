package internal

import (
	"fmt"
	"io"
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
