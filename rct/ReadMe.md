# Run CLI Tests

`rct` runs output-based tests for programs in a directory tree. From the
repository root, run it with:

```bash
go run ./rct [flags]
```

## Flags

* `--directory` (`-d`): Root directory to search. Defaults to the current
  directory (`.`).
* `--test-data-directory` (`-t`): Test data directory. Defaults to `./data`.
  This path is made absolute, and `TEST_DATA_DIR` in test commands is replaced
  with it.
* `--verbose` (`-v`): Print search, test setup, expected and actual output, and
  command exit codes.

## Test discovery and format

The CLI recursively searches the root and its subdirectories for directories
named `tests`. Hidden directories and directories whose names end in `.dSYM`
are skipped. In each `tests` directory, non-directory files ending in `.txt`
or `.err` (case-insensitive) are run in filename order; other files are
ignored.

Each test file must contain at least three lines:

1. A command to execute.
2. A separator line, which is ignored.
3. The expected output; any remaining lines are part of it.

For example:

```text
./my-program --option
---
Expected output
```

Commands are split on whitespace and run directly, not through a shell, with
the program directory as their working directory. The executable path is
resolved relative to that directory. `.txt` tests compare expected output with
standard output; `.err` tests compare with standard error. Comparison ignores
trailing newline characters, but otherwise requires an exact match.

The CLI prints a success count and percentage when tests are found. In verbose
mode it also displays each command's exit code. Exit codes do not affect
whether a test passes: pass/fail is determined only by output. The CLI exits
with status `1` when no tests are found or the operation is cancelled; test
failures alone do not change its exit status.
