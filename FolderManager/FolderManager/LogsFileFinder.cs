using NLog;

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

    public static Logger GetLogger(string tool, string task)
    {
        LogManager.Setup().LoadConfiguration(builder =>
        {
            var logFile = CreateLogsFile(tool, task);
            builder.ForLogger().WriteToFile(fileName: logFile);
        });

        var logger = LogManager.GetCurrentClassLogger();

        return logger;
    }
}
