# permreset

A simple command-line utility designed to find and correct file permissions for shell scripts within a specified directory.

This tool walks through a given directory tree, identifying files that appear to be executable scripts (i.e., those whose first line starts with `#!`, indicating a shebang). It then sets the file permissions for all found scripts to `0755` (rwxr-xr-x), ensuring they are executable by owner, group, and others.

## Prerequisites

*   Go installed on your system.
*   The target directory must exist and contain scripts with shebang lines.

## Usage

To run `permreset`, you must specify the root directory containing the scripts using the `--scriptsDirectory` (`-s`) flag.

```bash
# Syntax: go run main.go -s /path/to/scripts/root
go run . -s /path/to/your/local-scripts
```

**Example:**
If your scripts are located in `~/dev/scripts`, run:
```bash
go run . -s ~/dev/scripts
```

### Behavior

1.  **Script Detection:** Uses a check for `#!` on the first line to identify potential scripts.
2.  **Permission Check:** It will only modify permissions if the current mode is *not* already `0755`.
3.  **Platform Limitation:** This utility is designed for Unix-like systems (Linux/macOS) and will exit with an error if run on Windows.

## Running Tests

To ensure the tool works correctly, run the built-in tests:

```bash
go test ./...
```
