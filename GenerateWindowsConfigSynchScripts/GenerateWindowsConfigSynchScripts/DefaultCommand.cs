using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;

namespace GenerateWindowsConfigSynchScripts;

internal class DefaultCommand : AsyncCommand<CommandSettings>
{
    private readonly ILogger<WindowsConfigScriptsGenerator> _logger;

    public DefaultCommand(ILogger<WindowsConfigScriptsGenerator> logger)
    {
        ArgumentNullException.ThrowIfNull(logger);
        _logger = logger;
    }

    public override async Task<int> ExecuteAsync(CommandContext context, CommandSettings settings)
    {
        ArgumentNullException.ThrowIfNull(settings);

        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Script);

        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Autogen);
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Files);

        var synchFile = await SynchFileParser.ParseFile(settings.Files);

        var generator = new WindowsConfigScriptsGenerator(
            logger: _logger,
            autogen: settings.Autogen, 
            script: settings.Script,
            source: synchFile.Source,
            destination: synchFile.Destination,
            files: synchFile.Files);

        await generator.Generate();

        return 0;
    }
}
