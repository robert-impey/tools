````markdown
# staydeleted
To prevent files and directories from being accidentally deleted when using two-way synchronization tools (like rsync) with backup disks, `staydeleted` allows you to mark specific files or directories for protection.

This is particularly useful for folders synced to external backups on a regular schedule, where accidental deletions can be restored from the backup, bypassing the intended deletion cycle.

## ⚙️ How It Works: The Metadata File Mechanism

`staydeleted` does not modify the contents of your files. Instead, when you mark a file, it creates a hidden metadata file (e.g., `.stay-deleted/<hash>.txt`) inside a dedicated subdirectory, `.stay-deleted`, within the directory containing the protected file.

*   **Protection:** The metadata file contains instructions (`delete` or `keep`) that the `sweep` command reads.
*   **Dependency:** This behavior relies on your synchronization tools using options like `rsync --update` or similar mechanisms to read these metadata files and treat them as actual content.
*   **Cleanup:** The tool automatically manages the metadata files. Metadata files are marked for deletion and cleaned up after a set expiration period (default is 12 months), preventing the `.stay-deleted` folder from accumulating old, irrelevant data.

See [https://github.com/robert-impey/generate-synch-scripts](https://github.com/robert-impey/generate-synch-scripts) for related synchronization scripts.

---

### 🚀 Usage Guide

#### 1. Mark a File for Deletion or Keeping

Use `mark` to set the protection status for a file.

**Marking for Deletion (Default):**
This is the default action and marks the file for deletion, instructing the system to remove it.

```bash
PS C:\foo> staydeleted mark bar.txt
````

__Marking to Keep:__ Use the `--keep` flag (or `-k`) to tell the tool that this file must __not__ be deleted, even if it appears to be deleted from the primary source.

```bash
PS C:\foo> staydeleted mark --keep bar.txt
# OR
PS C:\foo> staydeleted mark -k bar.txt
```

#### 2. Sweep Directories (Standard Paths)

Use `sweep` to process directories that contain marked files.

__Syntax:__ `staydeleted sweep <dir1> [dir2] ...`

The tool will search each provided directory for metadata files, process the instructions, and perform the deletion of marked files.

```bash
# Sweep multiple directories at once
PS C:> staydeleted sweep C:\ProjectA C:\ProjectB
```

#### 3. Bulk Operations (From File)

For processing many files or directories, use the `*From` commands.

__Marking from a File (`markFrom`):__ Provide a text file containing one full file path per line. Lines starting with `#` are ignored, and empty lines are skipped.

```bash
# files-to-be-deleted.txt contains:
# # These files must be deleted
# C:\foo\old_document.pdf
# C:\foo\junk.txt
PS C:\foo> staydeleted markFrom files-to-be-deleted.txt
```

__Sweeping from a File (`sweepFrom`):__ Provide a text file containing one directory path per line. The tool will iterate through all listed directories and run `sweep` on each one.

```bash
# directories-to-sweep.txt contains:
# C:\ProjectA\Assets
# C:\ProjectB\Images
PS C:> staydeleted sweepFrom directories-to-sweep.txt
```

---

### 💡 Advanced Concepts

#### 📅 Expiry and Cleanup

The tool automatically manages the metadata files. By default, metadata files are marked for deletion and cleaned up after __12 months__. This prevents the `.stay-deleted` folder from accumulating old, irrelevant data.

You can explicitly override this cleanup period using the `--expiry` flag when running `sweep` or `sweepFrom`.

#### ⚠️ Handling Errors and Edge Cases

- __Multiple Paths:__ When multiple paths are provided (e.g., `sweep dir1 dir2`), the tool processes each path independently.
- __Missing Directories:__ If a directory provided to `sweep` or `sweepFrom` does not exist, the tool will log an error but will attempt to continue processing the remaining valid directories in the list.
- __Input File Format:__ For `markFrom` and `sweepFrom`, the tool supports comments (`#`) and skips empty lines.
- __Advanced Naming:__ The tool is designed to be robust and handles complex characters, spaces, and Unicode variations (e.g., NFC vs NFD) within file paths.
