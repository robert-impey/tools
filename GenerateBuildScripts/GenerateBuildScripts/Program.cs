using RunBuildScripts;

var folderManager = FolderManager.FolderManager.GetFolderManager();

var buildScriptFinder = new BuildScriptFinder(folderManager);

var buildScriptToCopy = buildScriptFinder.GetBuildScriptToCopy();

if (string.IsNullOrEmpty(buildScriptToCopy))
{
    Console.WriteLine("No build script found. Quitting...");
}
else
{
    Console.WriteLine($"Found {buildScriptToCopy}");

    var destination = buildScriptFinder.GetBuildScriptDestination();

    File.Copy(buildScriptToCopy, destination, overwrite: true);
}
