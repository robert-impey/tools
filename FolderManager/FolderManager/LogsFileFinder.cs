using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using NLog.Extensions.Logging;
using LogLevel = NLog.LogLevel;

namespace FolderManager;

public static class LogsFileFinder
{
    public static string CreateLogsFile(string tool, string task)
    {
        var synchLogsDirectory = Path.Combine(FolderManager.LogsFolder, tool);
        var timeString = DateTime.UtcNow.ToString("yyyy-MM-dd_HH-mm-ss");
        var fileName = $"{timeString}-{task}.log";
        return Path.Combine(synchLogsDirectory, fileName);
    }

    public static ILogger<T> GetLogger<T>(string tool, string task)
    {
        var logFile = CreateLogsFile(tool, task);
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
