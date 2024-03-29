namespace RunBuildScripts;

internal class BuildScriptRunner
{
    private readonly FolderManager.FolderManager _folderManager;

    public BuildScriptRunner(FolderManager.FolderManager folderManager)
    {
        _folderManager = folderManager;
    }

    public string? GetBuildScriptToRun()
    {
        var localScriptsFolder = _folderManager.GetLocalScriptsFolder();

        var machineLocalScriptsFolder = Path.Join(localScriptsFolder, Environment.MachineName);

        const string buildScript = "Build.ps1";

        var buildScriptToRun = Path.Join(machineLocalScriptsFolder, Environment.UserName, buildScript);

        if (File.Exists(buildScriptToRun))
        {
            return buildScriptToRun;
        }

        buildScriptToRun = Path.Join(machineLocalScriptsFolder, buildScript);

        if (File.Exists(buildScriptToRun))
        {
            return buildScriptToRun;
        }

        return null;
    }
}
