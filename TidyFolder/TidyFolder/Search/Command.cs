using Spectre.Console.Cli;

namespace TidyFolder.Search;

public sealed class Command : AsyncCommand<CommandSettings>
{
    protected override Task<int> ExecuteAsync(CommandContext context, CommandSettings settings, CancellationToken cancellationToken)
    {
        SearchLogic.SearchDirectory(settings.Directory, settings.LogsDirectory);

        return Task.FromResult(0);
    }
}
