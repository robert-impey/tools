﻿using FolderManager;
using Microsoft.Extensions.Logging;

namespace GenerateBuildScripts;

internal static class Program
{
    private static void Main(string[] args)
    {
        var logger = LogsFileFinder.GetLogger<FolderManager.FolderManager>("build", "GenerateBuildScripts") 
                     ?? throw new ApplicationException("Unable to create a logger!");

        var folderManager = FolderManager.FolderManager.GetFolderManager(logger);

        var buildScriptFinder = new BuildScriptFinder(folderManager);

        var buildScriptToCopy = buildScriptFinder.GetBuildScriptToCopy();

        var destination = buildScriptFinder.GetBuildScriptDestination();

        if (File.Exists(destination))
        {
            logger.LogInformation($"Deleting {destination}");
            File.Delete(destination);
        }

        if (string.IsNullOrEmpty(buildScriptToCopy))
        {
            logger.LogInformation("No build script found. Quitting...");
        }
        else
        {
            logger.LogInformation($"Found {buildScriptToCopy}");

            File.Copy(buildScriptToCopy, destination);
        }
    }
}