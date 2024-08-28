using GenerateBuildScripts;

var folderManager = FolderManager.FolderManager.GetFolderManager();

var buildScriptFinder = new BuildScriptFinder(folderManager);

var buildScriptToCopy = buildScriptFinder.GetBuildScriptToCopy();

var destination = buildScriptFinder.GetBuildScriptDestination();

if (File.Exists(destination))
{
    File.Delete(destination);
}

if (string.IsNullOrEmpty(buildScriptToCopy))
{
    Console.WriteLine("No build script found. Quitting...");
}
else
{
    Console.WriteLine($"Found {buildScriptToCopy}");

    File.Copy(buildScriptToCopy, destination);
}
