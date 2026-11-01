# LMF Tool Documentation

This repository contains the core logic for the LMF toolset, a sophisticated command-line utility designed to scan specified directory locations for a predefined set of managed folder names, and output a list of fully qualified paths.

## 📚 Project Overview

The tool is built using Go, utilizing the `cobra` framework for robust command-line argument parsing. Its primary function is to traverse multiple specified directories (`locationsFile`) and filter down a list of actual directories that match names provided in a separate file (`foldersFile`). The results are then saved to a designated output file.

## ⚙️ Code Structure & Components

The modular design ensures that core logic is separated from the command handling:

*   **`main.go` (cmd package)**:
    *   **Role**: The command entry point. It initializes the CLI structure and orchestrates the execution flow.
    *   **Key Functionality**: It handles input file validation (locations, folders) and manages the entire I/O process. It calls `common.LoadFolderManager` to ingest paths and then passes the resulting manager object to the `internal.FolderManager` for listing.
    *   **Utilities**: Contains `writeAutoGenHeader`, which stamps the output file with a mandatory, non-editable metadata header (including creation date).

*   **`internal/` Package (Core Logic)**: This package holds the core business logic.
    *   **`list_managed_folders.go`**: Defines the `FolderManager` struct and the primary methods.
        *   `ListManagedFolders(outFile io.Writer)`: This method iterates through all defined locations. For each location, it uses `os.Lstat` to verify if the path points to an actual directory (and not a symbolic link or file). It ensures that location outputs are cleanly separated by blank lines.
        *   `GetCleanLocationName`: Performs essential path sanitization by replacing common path separators (`:`, `\`, `/`, ` `) with underscores (`_`) and trimming any resulting trailing underscores.
        *   `GetScriptPath`: A utility to correctly join directory paths using `filepath.Join`, mimicking Java's `Path.resolve()`.

*   **`cmd/` Package**:
    *   **`root.go`**: Defines the main `rootCmd` structure, setting the CLI usage (`foldermanager`) and defining flags (`-l`/`--locations`, `-f`/`--folders`, `-m`/`--managed-folders-file`).

## 🚀 How to Run (Execution Flow)

The application requires three pieces of configuration data:

1.  **Locations File (`-l` or `--locations`)**: A file containing absolute paths to the directories to scan (one path per line).
2.  **Folders File (`-f` or `--folders`)**: A file containing the exact names of the folders you are looking for (one name per line).
3.  **Output File (`-m` or `--managed-folders-file`)**: The path where the final list of absolute paths will be written.

**Execution Example:**
```bash
./lmf --locations /path/to/locs --folders /path/to/folders --managed-folders-file /output/lmf_output.txt
```

## 🧪 Testing and Validation

The core logic is thoroughly tested in `internal/list_managed_folders_test.go`. These tests validate:
*   Correct path sanitization.
*   Accurate listing of managed folders while ignoring non-directory files or symbolic links.
*   Proper formatting, including inserting blank lines between separate location output blocks.

## 🤝 Contributing Guidelines

Please ensure all new features adhere to the existing structural patterns. When modifying the core logic in `internal/`, always utilize the provided tests to ensure backward compatibility and correctness.
