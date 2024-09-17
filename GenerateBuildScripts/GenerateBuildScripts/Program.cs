using FolderManager;
using GenerateBuildScripts;

internal class Program
{
    private static void Main(string[] args)
    {
        var logger = LogsFileFinder.GetLogger("build", "GenerateBuildScripts") ?? throw new ApplicationException("Unable to create a logger!");

        var folderManager = FolderManager.FolderManager.GetFolderManager(logger);

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