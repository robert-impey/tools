using FolderManager;
using GenerateBuildScripts;
using NLog;

internal class Program
{
    private static void Main(string[] args)
    {
        LogManager.Setup().LoadConfiguration(builder =>
        {
            var logFile = LogsFileFinder.CreateLogsFile("build", "GenerateBuildScripts");
            builder.ForLogger().WriteToFile(fileName: logFile);
        });

        var logger = LogManager.GetCurrentClassLogger();

        var folderManager = FolderManager.FolderManager.GetFolderManager(logger as Microsoft.Extensions.Logging.ILogger);

        var buildScriptFinder = new BuildScriptFinder(folderManager);

        var buildScriptToCopy = buildScriptFinder.GetBuildScriptToCopy();

        var destination = buildScriptFinder.GetBuildScriptDestination();


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