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

    var p = new Process();
    p.StartInfo.FileName = folderManager.PowerShellExe;
    p.StartInfo.ArgumentList.Add(buildScriptToRun);
    p.Start();
}
