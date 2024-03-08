var folderManager = FolderManager.FolderManager.GetFolderManager();

var buildScriptToRun = folderManager.GetBuildScriptToRun();

if (string.IsNullOrEmpty(buildScriptToRun))
{
    Console.WriteLine("No build script to run. Quitting...");
}
else
{
    Console.WriteLine($"Running {buildScriptToRun}");
}
