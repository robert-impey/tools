﻿using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;

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

        using var servicesProvider = new ServiceCollection()
            .AddLogging(loggingBuilder =>
            {
                loggingBuilder.ClearProviders();
                loggingBuilder.AddConsole();
            }).BuildServiceProvider();
        
        var logger = servicesProvider.GetRequiredService<ILogger<T>>();
        
        return logger;
    }
}