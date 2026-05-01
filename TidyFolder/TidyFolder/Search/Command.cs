using Spectre.Console.Cli;

namespace TidyFolder.Search;

public sealed class Command : AsyncCommand<CommandSettings>
{
    protected override Task<int> ExecuteAsync(CommandContext context, CommandSettings settings, CancellationToken cancellationToken)
    {
        Console.WriteLine($"{settings.Directory}");

        if (string.IsNullOrEmpty(settings.LogsDirectory))
        {
            Console.WriteLine("Logs directory is not set.");
        }
        else
        {
            Console.WriteLine($"Logs directory: {settings.LogsDirectory}");
        }

        return Task.FromResult(0);
    }
}
