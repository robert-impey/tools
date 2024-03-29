using RunBuildScripts;
using System.Diagnostics;

var folderManager = FolderManager.FolderManager.GetFolderManager();

var buildScriptRunner = new BuildScriptRunner(folderManager);

var buildScriptToRun = buildScriptRunner.GetBuildScriptToRun();

if (string.IsNullOrEmpty(buildScriptToRun))
{
    Console.WriteLine("No build script to run. Quitting...");
}
else
{
    Console.WriteLine($"Running {buildScriptToRun}");

    Process p = new Process();
    p.StartInfo.FileName = buildScriptToRun;
    p.Start();
}
