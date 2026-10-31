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
/cygdrive/x
/cygdrive/c/Users/rober/OneDrive

config
data
docs
local-scripts
```

Save this as `OneDrive.txt`, then run `gensyn` from PowerShell, with the output
directory already created:

```powershell
gensyn.exe rsync --autogenDir "$HOME\autogen\synch" .\OneDrive.txt
```

The input file has four parts: the `rsync` command and its options, the source
path, the destination path, a blank line, and then the names of the directories
to synchronize. The input filename determines the generated script names.

In the default directory mode, this example creates
`$HOME\autogen\synch\OneDrive.ps1`. Running it executes the specified `rsync`
command for each listed directory, first from source to destination and then
in the reverse direction. Since the example lists multiple directories,
`gensyn` also creates `OneDrive\config.ps1`, `OneDrive\data.ps1`,
`OneDrive\docs.ps1`, and `OneDrive\local-scripts.ps1`; each runs the same
two-way synchronization for just that directory. The generated scripts load
SSH-agent settings from the user's keychain/environment files when present.

Pass `--files` (or `-f`) to generate file-mode commands instead. In this mode,
each listed name is treated as a file appended to the source path, and the
destination is used as given. Only the combined `<name>.ps1` script is
generated. Generated scripts accept the `-logged` switch; when used, they
write each direction's output and errors under the user's `logs\synch`
directory.
