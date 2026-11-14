package cronsetter

/*
Copyright © 2023 Robert Impey robert-impey@users.noreply.github.com
*/

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type refreshState int

const (
	stateBeforeHeader refreshState = iota
	stateInsideManagedBlock
	stateAfterFooter
)

func RefreshCrontab(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	state := stateBeforeHeader

	for scanner.Scan() {
		line := scanner.Text()

		switch state {
		case stateBeforeHeader:
			if isHeaderComment(line) {
				PrintTasks(w)
				state = stateInsideManagedBlock
			} else {
				fmt.Fprintln(w, line)
			}
		case stateInsideManagedBlock:
			if isFooterComment(line) {
				state = stateAfterFooter
			}
		case stateAfterFooter:
			fmt.Fprintln(w, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if state == stateBeforeHeader {
		PrintTasks(w)
	}

	return nil
}

func isHeaderComment(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), headerCommentPrefix)
}

func isFooterComment(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), footerCommentPrefix)
}
