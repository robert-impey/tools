using FolderManager;
using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;

namespace GenerateWindowsConfigSynchScripts;

public class Command : AsyncCommand<CommandSettings>
{
    public override async Task<int> ExecuteAsync(CommandContext context, CommandSettings settings)
    {
        ArgumentNullException.ThrowIfNull(settings);

        ArgumentException.ThrowIfNullOrWhiteSpace(settings.LogsDirectory);
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Script);
        
        ILogger<WindowsConfigScriptsGenerator> logger;

        if (settings.Logged)
        {
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
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Source);
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Destination);

        var generator = new WindowsConfigScriptsGenerator(
            logger: logger,
            logsDirectory: settings.LogsDirectory,
            autogen: settings.Autogen, 
            script: settings.Script,
            files: settings.Files,
            source: settings.Source,
            destination: settings.Destination);

        await generator.Generate();

        return 0;
    }
}
