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
}
