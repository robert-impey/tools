# Dupes

A fast, cross-platform command-line utility written in Dart to find duplicate files across one or more directories. 

It optimizes the search by first grouping files by exact byte size, and only calculating SHA256 checksums for files that share the exact same size, saving processing time.

## Installation

### Prerequisites

To use Dupes, you must have the Dart SDK installed on your system. You can download it from the official [Dart website](https://dart.dev/get-dart).

### Getting Started

1. Clone or download this repository to your local machine.
2. Open a terminal and navigate to the root directory of the project.
3. Fetch the required dependencies:

```bash
dart pub get
```

## Usage

Run the tool by passing one or more directories as arguments. It will search files recursively through the directories.

```bash
dart run bin/dupes.dart <directory1> [directory2] ...
```

### Output modes

By default, the tool prints a human-readable list of duplicate files grouped by their SHA-256 checksum.

```bash
dart run bin/dupes.dart /path/to/photos
```
Example Output:
```
--- Checksum: a1b2c3d4e5f6... ---
    /path/to/photos/vacation/image1.jpg
    /path/to/photos/backup/image1_copy.jpg

--- Checksum: f6e5d4c3b2a1... ---
    /path/to/photos/doc.pdf
    /path/to/photos/archive/doc.pdf
```

### `--dupes` flag

If you want to automate the cleanup process and pipe the output to other tools (like `rm`), use the `--dupes` flag. This mode suppresses the headers and checksums. For every group of identical files, it preserves the first file found (alphabetically) and prints only the paths of the *subsequent* duplicates, one per line.

```bash
dart run bin/dupes.dart --dupes /path/to/photos
```

Example Output:
```
/path/to/photos/backup/image1_copy.jpg
/path/to/photos/archive/doc.pdf
```

Delete all duplicates using xargs (Linux/macOS):
```bash
dart run bin/dupes.dart --dupes /path/to/photos | xargs rm
```