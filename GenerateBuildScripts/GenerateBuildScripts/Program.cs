using GenerateBuildScripts;
using NLog;

internal class Program
{
    private static void Main(string[] args)
    {
        var folderManager = FolderManager.FolderManager.GetFolderManager();

        var buildScriptFinder = new BuildScriptFinder(folderManager);

        var buildScriptToCopy = buildScriptFinder.GetBuildScriptToCopy();

        var destination = buildScriptFinder.GetBuildScriptDestination();

        LogManager.Setup().LoadConfiguration(builder =>
        {
            var synchLogsDirectory = Path.Combine(folderManager.GetLogsFolder(), "build");
            var timeString = DateTime.UtcNow.ToString("yyyy-MM-dd_HH-mm-ss");
            var fileName = $"{timeString}-GenerateBuildScripts.log";
            var logFile = Path.Combine(synchLogsDirectory, fileName);
            builder.ForLogger().WriteToFile(fileName: logFile);
        });

        var logger = LogManager.GetCurrentClassLogger();

        if (File.Exists(destination))
        {
            logger.Info($"Deleting {destination}");
            File.Delete(destination);
        }

        if (string.IsNullOrEmpty(buildScriptToCopy))
        {
            logger.Info("No build script found. Quitting...");
        }
        else
        {
            logger.Info($"Found {buildScriptToCopy}");

            File.Copy(buildScriptToCopy, destination);
        }
    }
}