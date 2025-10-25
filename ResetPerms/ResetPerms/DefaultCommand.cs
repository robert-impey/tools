using FolderManager;
using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;
using System.Collections.Immutable;

namespace ResetPerms;

public class DefaultCommand : Command<CommandSettings>
{
    public override int Execute(CommandContext context, CommandSettings settings)
    {
        ILogger logger;

        if (settings.Logged)
        {
            if (string.IsNullOrWhiteSpace(settings.LogsDirectory))
            {
                throw new ArgumentNullException(nameof(settings.LogsDirectory));
            }

            logger = LogsFileFinder.GetLogger<DefaultCommand>(settings.LogsDirectory, "ResetPerms");
        }
        else
        {
            using var loggerFactory = LoggerFactory.Create(builder =>
            {
                builder.AddConsole();
            });
            logger = loggerFactory.CreateLogger<DefaultCommand>();
        }

        var filesWithShebang = ScriptsFinder.FindFilesWithShebang(settings.ScriptsDirectory).ToImmutableArray();

        logger.LogInformation("Found {Length} files with shebangs", filesWithShebang.Length);

        foreach (var file in filesWithShebang)
        {
            logger.LogInformation("File with shebang: {File}", file);

#pragma warning disable CA1416
            File.SetUnixFileMode(
                file,
                UnixFileMode.UserRead
                | UnixFileMode.UserWrite
                | UnixFileMode.UserExecute
                | UnixFileMode.GroupRead
                | UnixFileMode.GroupExecute
                | UnixFileMode.OtherRead
                | UnixFileMode.OtherExecute
            );
#pragma warning restore CA1416
        }

        return 0;
    }
}
