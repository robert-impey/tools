# Run CLI Tests


## Documentation

This repository contains the source code and tooling for the `rct` CLI, which is designed to run and validate tests across multiple program directories.

### Usage

The core functionality is exposed through the `main.go` file, which handles command-line argument parsing and orchestrates the test execution process.

The CLI is invoked using:

```bash
./main.go [flags]
```

### Flags

*   `--directory` (`-d`): The root directory to search for programs. Defaults to the current directory (`.`).
*   `--test-data-directory` (`-t`): Directory containing test data. Defaults to `./data`.
*   `--verbose` (`-v`): Enables verbose output, showing detailed test setup and execution logs.

### File Structure

*   **`main.go`**: Contains the main entry point, argument parsing logic, and the core `execute` function responsible for traversing the file system and running tests.
*   **`ReadMe.md`**: This file provides an overview of how to use the CLI and the project structure.

***

