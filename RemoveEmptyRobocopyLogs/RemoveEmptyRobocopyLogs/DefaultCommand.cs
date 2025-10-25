using Microsoft.Extensions.FileSystemGlobbing;
using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;

namespace RemoveEmptyRobocopyLogs;

public sealed class DefaultCommand : Command<CliSettings>
{
    // C# uses constructor injection for ILogger, which is set up in Program.cs
    private readonly ILogger<DefaultCommand> _logger;

    public DefaultCommand(ILogger<DefaultCommand> logger)
    {
        _logger = logger;
    }

    public override int Execute(CommandContext context, CliSettings settings)
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
            // If the directory doesn't exist, we should probably exit with a non-zero code.
            return 1;
        }

        // Setup Globbing Matcher
        var matcher = new Matcher();
        // The F# code used a sequence; in C#, we add patterns directly.
        matcher.AddInclude("**/*.robocopy-synch.log"); // Using ** for recursive search is common, or just * if non-recursive

        // Get matching files
        // Note: GetResultsInFullPath returns a string enumerable directly in C#
        var matchingFiles = matcher.GetResultsInFullPath(settings.LogsDirectory).ToList();

        _logger.LogInformation("There are {FileCount} log files", matchingFiles.Count);

        // --- File Processing Loop ---
        foreach (var logFile in matchingFiles)
        {
            try
            {
                // This assumes RobocopyLogs.FileHasCopies is a static method as per previous translation
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

        return 0; // Success
    }
}
