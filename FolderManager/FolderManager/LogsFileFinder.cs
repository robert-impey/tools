using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using NLog.Extensions.Logging;
using LogLevel = NLog.LogLevel;

namespace FolderManager;

public static class LogsFileFinder
{
    public const string LogFileTimeFormat = "yyyy-MM-ddTHH-mm-ssZ";

    public static string CreateLogsFile(string logsDirectory, string tool, string task)
    {
        var toolLogsDirectory = Path.Combine(logsDirectory, tool);
        var timeString = DateTime.UtcNow.ToString(LogFileTimeFormat);
        var fileName = $"{timeString}-{task}.log";
        return Path.Combine(toolLogsDirectory, fileName);
    }

    public static ILogger<T> GetLogger<T>(string logsDirectory, string tool, string task)
    {
        var logFile = CreateLogsFile(logsDirectory, tool, task);
        var config = new NLog.Config.LoggingConfiguration();
        
        var logFileTarget = new NLog.Targets.FileTarget("logfile") { FileName = logFile };
        
        config.AddRule(LogLevel.Trace, LogLevel.Fatal, logFileTarget);
        
        using var servicesProvider = new ServiceCollection()
            .AddLogging(loggingBuilder =>
            {
                loggingBuilder.ClearProviders();
                loggingBuilder.AddNLog(config);
            }).BuildServiceProvider();
        
        var logger = servicesProvider.GetRequiredService<ILogger<T>>();
        
        return logger;
    }
}
