using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;
using System.Collections.Immutable;

namespace ResetPerms;

public class DefaultCommand : AsyncCommand<CommandSettings>
{
    private readonly ILogger<DefaultCommand> _logger;

    public DefaultCommand(ILogger<DefaultCommand> logger)
    {
        ArgumentNullException.ThrowIfNull(logger);
        _logger = logger;
    }

    public override Task<int> ExecuteAsync(CommandContext context, CommandSettings settings, CancellationToken cancellationToken)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.ScriptsDirectory);

        var filesWithShebang = ScriptsFinder.FindFilesWithShebang(settings.ScriptsDirectory).ToImmutableArray();

        _logger.LogInformation("Found {Length} files with shebangs", filesWithShebang.Length);

        foreach (var file in filesWithShebang)
        {
            _logger.LogInformation("File with shebang: {File}", file);

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

        return Task.FromResult(0);
    }
}
