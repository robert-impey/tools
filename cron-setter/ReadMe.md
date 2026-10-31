# cron-setter

`cron-setter` generates crontab entries to schedule the daily maintenance scripts (such as `stay_deleted`, directory synchronisation, permission resets, builds, log cleanup, and managed folder listing) at randomised minutes within their designated hours.

The generated entries are wrapped in header and footer comments so they can be identified and updated without overwriting custom cron jobs:

```text
# Start of crontab tasks created by cron-setter at <YYYY-MM-DD HH:MM:SS>
...
# End of cron tasks created by cron-setter at <YYYY-MM-DD HH:MM:SS>
```

## Usage

### Create a new crontab (`print`)

`cron-setter print` outputs a complete set of `cron-setter` tasks to `stdout`. To replace your current crontab with the generated tasks and inspect the result:

```sh
cron-setter print | crontab - && crontab -l
```

### Refresh the managed tasks while keeping custom entries (`refresh`)

`cron-setter refresh` reads an existing crontab from `stdin`, preserves any lines before the `# Start of crontab tasks created by cron-setter` header and after the `# End of cron tasks created by cron-setter` footer, and regenerates the `cron-setter` tasks in between.

If no `cron-setter` section exists yet in the input, the generated block is appended to the end.

To refresh the `cron-setter` section of your current crontab in place:

```sh
crontab -l | cron-setter refresh | crontab - && crontab -l
```

For example, given a crontab with custom tasks before and after the managed block:

```crontab
# Custom task before cron-setter
0 4 * * * /path/to/morning-job.sh

# Start of crontab tasks created by cron-setter at 2026-10-31 10:58:48
...
# End of cron tasks created by cron-setter at 2026-10-31 10:58:48

# Custom task after cron-setter
30 23 * * * /path/to/evening-job.sh
```

Running `crontab -l | cron-setter refresh | crontab -` keeps both custom tasks intact and generates fresh schedules for the tasks inside the `cron-setter` block.

## Configuration

By default, `cron-setter` resolves script paths relative to `~/local-scripts`. You can override this base directory by setting the `LOCAL_SCRIPTS` environment variable.
