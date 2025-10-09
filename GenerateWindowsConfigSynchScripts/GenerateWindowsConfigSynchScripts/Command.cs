using FolderManager;
using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;

namespace GenerateWindowsConfigSynchScripts;

public class Command : AsyncCommand<CommandSettings>
{
    public override async Task<int> ExecuteAsync(CommandContext context, CommandSettings settings)
    {
        ArgumentNullException.ThrowIfNull(settings);

        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Script);
        
        ILogger<WindowsConfigScriptsGenerator> logger;

        if (settings.Logged)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(settings.LogsDirectory);
            logger = LogsFileFinder.GetLogger<WindowsConfigScriptsGenerator>(
                settings.LogsDirectory, $"WindowsConfigScriptsGenerator-{settings.Script}");
        }
        else
        {
            using var loggerFactory = LoggerFactory.Create(builder =>
            {
                builder.AddConsole();
            });
            logger = loggerFactory.CreateLogger<WindowsConfigScriptsGenerator>();
        }

        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Autogen);
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Files);

        var synchFile = await SynchFileParser.ParseFile(settings.Files);

        var generator = new WindowsConfigScriptsGenerator(
            logger: logger,
            autogen: settings.Autogen, 
            script: settings.Script,
            source: synchFile.Source,
            destination: synchFile.Destination,
            files: synchFile.Files);

        await generator.Generate();

        return 0;
    }
}
