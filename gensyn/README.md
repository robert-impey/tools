# gensyn

`gensyn` generates PowerShell scripts for synchronizing folders and files.

These scripts can be invoked as scheduled tasks.

Note that synchronizing two folders in this way can make file deletion a problem.
This tool may help: https://github.com/robert-impey/tools/tree/prod/staydeleted

## Commands

`gensyn` has three subcommands:

It currently supports:

| Command | Purpose |
| --- | --- |
| `rsync` | Generate PowerShell scripts that run `rsync` in both directions. |
| `robocopy` | Generate PowerShell scripts that run `ROBOCOPY` for managed folder sets. |
| `rcif` | Generate a PowerShell script that runs `ROBOCOPY` in both directions for individual files. |

Use the command-specific sections below for the required flags and input file formats.

### Rsync

```
rsync --update --recursive --verbose --times --iconv=utf8 
/mnt/x
/home/robert/Dropbox

config
data
docs
local-scripts
```

Save this as `Dropbox.txt`, then run `gensyn` from a Linux or macOS shell,
with the output directory already created:

```sh
gensyn rsync --autogenDir "$HOME/autogen/synch" ./Dropbox.txt
```

The input file has four parts: the `rsync` command and its options, the source
path, the destination path, a blank line, and then the names of the directories
to synchronize. The input filename determines the generated script names.

In the default directory mode, this example creates
`$HOME/autogen/synch/Dropbox.ps1`. Running it executes the specified `rsync`
command for each listed directory, first from source to destination and then
in the reverse direction. Since the example lists multiple directories,
`gensyn` also creates `Dropbox/config.ps1`, `Dropbox/data.ps1`,
`Dropbox/docs.ps1`, and `Dropbox/local-scripts.ps1`; each runs the same
two-way synchronization for just that directory. The generated scripts load
SSH-agent settings from the user's keychain/environment files when present.

Pass `--files` (or `-f`) to generate file-mode commands instead. In this mode,
each listed name is treated as a file appended to the source path, and the
destination is used as given. Only the combined `<name>.ps1` script is
generated. Generated scripts accept the `-logged` switch; when used, they
write each direction's output and errors under the user's `logs/synch`
directory.

### Robocopy

Create a locations file, for example `Locations.txt`:

```
C:\Data
D:\Backup
```

Create a folders file, for example `Folders.txt`:

```
Documents
Pictures
```

Each nonblank, non-comment line in the locations file is a location root;
each line in the folders file is a managed folder name relative to those
roots. Lines beginning with `#` are comments. Then run `gensyn` from
PowerShell, with the output directory already created:

```powershell
gensyn robocopy --locations .\Locations.txt --folders .\Folders.txt --autogen C:\Users\Robert\autogen\synch
```

The flags also have short forms: `-l` for `--locations`, `-f` for
`--folders`, and `-a` for `--autogen`. For each ordered pair of distinct
locations, `gensyn` creates a subdirectory named after the two locations
(with path separators, colons, and spaces replaced by underscores). For the
example above, it generates scripts under both
`C_Data\D_Backup` and `D_Backup\C_Data`. Each pair gets a script for every
listed folder present in both locations, as well as an `_all.ps1` script
that runs all of those folders. The generated scripts copy in one direction
only, so run the script from the other location pair for the reverse
direction.

The scripts invoke `ROBOCOPY` with `/E /XO /R:0`: they include subdirectories
(including empty ones), skip older source files, and do not retry failed
copies. A script checks that each folder exists at both locations before
copying. Pass the `-logged` switch when running a generated script to write
its output and errors to `.log` and `.err` files under the user's
`logs/synch` directory.
