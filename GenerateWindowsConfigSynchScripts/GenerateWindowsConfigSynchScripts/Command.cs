using FolderManager;
using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;

namespace GenerateWindowsConfigSynchScripts;
public class Command : Command<CommandSettings>
{
    public override int Execute(CommandContext context, CommandSettings settings)
    {
        ArgumentNullException.ThrowIfNull(settings);

        ArgumentNullException.ThrowIfNull(settings.LogsDirectory);

        ILogger logger;

        if (settings.Logged)
        {
            logger = LogsFileFinder.GetLogger<Command>(settings.LogsDirectory, "synch", "GenerateWindowsConfigSynchScripts");
        }
        else
        {
            using var loggerFactory = LoggerFactory.Create(builder =>
            {
                builder.AddConsole();
            });
            logger = loggerFactory.CreateLogger<Command>();
        }
        
        logger.LogInformation("Starting GenerateWindowsConfigSynchScripts");

        if (string.IsNullOrWhiteSpace(settings.Autogen))
        {
            logger.LogError("No autogen directory specified, exiting...");
            return 1;
        }
        logger.LogInformation($"Autogen: {settings.Autogen}");

        if (string.IsNullOrWhiteSpace(settings.Files))
        {
            logger.LogError("No file listing files specified, exiting...");
            return 1;
        }
        logger.LogInformation($"Files: {settings.Files}");

        if (string.IsNullOrWhiteSpace(settings.Source))
        {
            logger.LogError("No source folder specified, exiting...");
            return 1;
        }
        logger.LogInformation($"Source: {settings.Source}");

        if (string.IsNullOrWhiteSpace(settings.Destination))
        {
            logger.LogError("No destination folder specified, exiting...");
            return 1;
        }
        logger.LogInformation($"Destination: {settings.Destination}");


        return 0;
    }
}
