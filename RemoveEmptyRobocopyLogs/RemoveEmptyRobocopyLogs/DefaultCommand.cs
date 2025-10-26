using Microsoft.Extensions.FileSystemGlobbing;
using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;

namespace RemoveEmptyRobocopyLogs;

public sealed class DefaultCommand : Command<CommandSettings>
{
    private readonly ILogger<DefaultCommand> _logger;

    public DefaultCommand(ILogger<DefaultCommand> logger)
    {
        _logger = logger;
    }

    public override int Execute(CommandContext context, CommandSettings settings)
    {
        _logger.LogInformation("Looking for Robocopy Log Files");
        _logger.LogInformation("Logs directory: {LogsDirectory}", settings.LogsDirectory);

        if (Directory.Exists(settings.LogsDirectory))
        {
            _logger.LogInformation("Synch logs directory {LogsDirectory} exists", settings.LogsDirectory);
        }
        else
        {
            _logger.LogWarning("Synch logs directory {LogsDirectory} does not exist. Exiting.", settings.LogsDirectory);
            return 1;
        }

        var matcher = new Matcher();
        matcher.AddInclude("**/*.robocopy-synch.log");

        var matchingFiles = matcher.GetResultsInFullPath(settings.LogsDirectory).ToList();

        _logger.LogInformation("There are {FileCount} log files", matchingFiles.Count);

        foreach (var logFile in matchingFiles)
        {
            try
            {
                if (RobocopyLogsParser.FileHasCopies(logFile))
                {
                    _logger.LogInformation("{LogFile} has copies - keeping", logFile);
                }
                else
                {
                    _logger.LogInformation("{LogFile} has no copies - deleting", logFile);
                    File.Delete(logFile);
                }
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to process or delete file: {LogFile}", logFile);
            }
        }

        return 0; 
    }
}
