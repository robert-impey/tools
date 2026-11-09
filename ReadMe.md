# Tools

The programs in this repository are tools that help me curate a large number of files in my home network.
I don't like repetitive tasks, but I do find it fun to automate stuff.

I've used a number of different programming languages in order to learn those languages.
Learning new languages by developing and improving programs is also an aim of this project.
While multi-language development is considerably less efficient from the point of view of software
engineering than writing everything in one language, I have benefitted from seeing how different
languages solve the same problems.
The number of languages in use has gone up and down over time.
I've found that in order for knowledge to stick, I need to solve problems that have some utility.
This is in contrast to the toy programs that I have written at https://github.com/robert-impey/CodingExperiments

## Tool Summaries

Each tool is designed to solve a practical automation problem:

- **gensyn**: Generates PowerShell scripts for synchronizing folders/files via `rsync`, `ROBOCOPY`, or per-file sync commands; intended for scheduled task execution.
- **sorttrees**: Sorts space-indented text trees by grouping items at the same indentation level, useful for organizing hierarchical lists like to-do items.
- **staydeleted**: Manages coordinated multi-site file deletion across synchronized directories using `.stay-deleted` metadata files with automatic 12-month expiry cleanup.
- **permreset**: Fixes executable script permissions on Unix-like systems by identifying shebang (`#!`) scripts and setting their mode to `0755`.
- **dupes**: Fast cross-platform duplicate file finder that groups files by size first, then computes SHA256 checksums only for same-sized candidates.
- **lmf**: Lists managed folders by scanning specified directories against a list of folder names and outputting matching fully qualified paths.
- **logs-deleter**: Cleans up old logs (rsync, robocopy, custom) using retention policies based on file age and content analysis; supports per-tool sweeps.
- **rct**: Output-based test runner that discovers `.txt`/`.err` test files in `tests/` directories and compares command outputs against expected results.
- **cron-setter**: Generates crontab entries to schedule maintenance scripts at randomised minutes within designated hours, wrapped in comment headers to preserve custom jobs.
- **TidyFolder (search)**: Scans a directory tree for related files with matching extensions where one filename stem starts with another, grouped by parent folder.