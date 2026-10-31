# TidyFolder

## `search` command

The `search` command scans a directory tree and reports files in the same folder that have matching extensions and one filename stem starts with the other. 
It can be used to find related files such as `report.txt` and `report-final.txt`.

```text
TidyFolder search [OPTIONS] <directory>
```

### Arguments and options

- `<directory>` — required path to an existing directory. The command searches its files recursively.
- `-l|--logs-dir <PATH>` — optional path to an existing directory where output and error log files are written.

For example:

```text
TidyFolder search "C:\\Documents"
TidyFolder search --logs-dir "C:\\Logs" "C:\\Documents"
```

### Matching and output

Only files are considered; directories are not included.
A pair is reported when both files are in the same folder, their extensions match without regard to case, and one filename stem starts with the other using a case-sensitive comparison. 
For example, `report.txt` and `report-final.TXT` match, but `Report-final.txt` does not match `report.txt` because the stem comparison is case-sensitive.
A file with no extension can match another extensionless file.

Results are grouped by parent folder, with the two filenames listed under each folder. 
If there are no matches, the command prints no results.

If matches are found and `--logs-dir` is supplied, results are written to a UTF-8 log named `<timestamp>-search-<directory>.log`; an error log with the same timestamp and `.err` extension is also created. 
The directory portion of these filenames is sanitized for use as a filename.
If a log file cannot be created, a warning is written to standard error. Without this option, results are written to standard output.
