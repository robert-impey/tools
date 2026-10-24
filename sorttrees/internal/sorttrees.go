package internal

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
)

const indentationWidth = 4

type parsedLine struct {
	depth int
	text  string
}

type tree struct {
	text     string
	children []*tree
}

func ReadLines(fileName string) ([]string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	rawLines := strings.Split(string(data), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		line = strings.TrimRightFunc(line, unicode.IsSpace)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func SortLines(lines []string, separateTopLevel bool) (string, error) {
	root := &tree{}
	stack := []*tree{root}

	for lineNumber, line := range lines {
		parsed := parseLine(line)
		if parsed.depth >= len(stack) {
			return "", fmt.Errorf("line %d skips indentation depth", lineNumber+1)
		}

		parent := stack[parsed.depth]
		child := &tree{text: parsed.text}
		parent.children = append(parent.children, child)

		stack = stack[:parsed.depth+1]
		stack = append(stack, child)
	}

	finalise(root)
	var output strings.Builder
	writeTree(&output, root, "", separateTopLevel, true)
	return output.String(), nil
}

func AreLinesSortedTree(lines []string) (bool, error) {
	previousItemsByDepth := make([]string, 0)
	previousLineDepth := 0

	for lineNumber, line := range lines {
		parsed := parseLine(line)

		if parsed.depth < previousLineDepth {
			oldItems := previousItemsByDepth
			if parsed.depth >= len(oldItems) {
				return false, fmt.Errorf("line %d has invalid indentation depth", lineNumber+1)
			}
			previousItemsByDepth = append([]string(nil), oldItems[:parsed.depth+1]...)
		}

		if parsed.depth == len(previousItemsByDepth) {
			previousItemsByDepth = append(previousItemsByDepth, parsed.text)
		} else if parsed.depth < len(previousItemsByDepth) {
			if parsed.text < previousItemsByDepth[parsed.depth] {
				return false, nil
			}
			previousItemsByDepth[parsed.depth] = parsed.text
		} else {
			return false, fmt.Errorf("line %d skips indentation depth", lineNumber+1)
		}

		previousLineDepth = parsed.depth
	}

	return true, nil
}

func BackupFile(fileName string) error {
	info, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	source, err := os.Open(fileName)
	if err != nil {
		return err
	}

	backupName := fileName + ".bak"
	backup, err := os.OpenFile(backupName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		source.Close()
		return err
	}
	_, copyErr := io.Copy(backup, source)
	closeErr := backup.Close()
	sourceCloseErr := source.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if sourceCloseErr != nil {
		return sourceCloseErr
	}
	if err := os.Chmod(backupName, info.Mode()); err != nil {
		return err
	}
	return os.Chtimes(backupName, info.ModTime(), info.ModTime())
}

func parseLine(line string) parsedLine {
	indentation := 0
	for indentation < len(line) && line[indentation] == ' ' {
		indentation++
	}
	return parsedLine{
		depth: indentation / indentationWidth,
		text:  strings.TrimRightFunc(line[indentation:], unicode.IsSpace),
	}
}

func finalise(node *tree) {
	for _, child := range node.children {
		finalise(child)
	}
	sortTrees(node.children)
}

func sortTrees(trees []*tree) {
	sort.SliceStable(trees, func(i, j int) bool {
		return lessTree(trees[i], trees[j])
	})
}

func lessTree(left, right *tree) bool {
	if left.text < right.text {
		return true
	}
	return treeString(left) < treeString(right)
}

func treeString(node *tree) string {
	var result strings.Builder
	writeTree(&result, node, "", false, false)
	return result.String()
}

func writeTree(output *strings.Builder, node *tree, indentation string, separateTopLevel, isRoot bool) {
	if !isRoot {
		output.WriteString(indentation)
		output.WriteString(node.text)
		output.WriteByte('\n')
		indentation += strings.Repeat(" ", indentationWidth)
	}

	for i, child := range node.children {
		if isRoot && separateTopLevel && i > 0 {
			output.WriteByte('\n')
		}
		writeTree(output, child, indentation, separateTopLevel, false)
	}
}
