namespace RunBuildScripts;

internal class BuildScriptFinder
{
    private readonly FolderManager.FolderManager _folderManager;

    private const string buildScript = "Build.ps1";

    public BuildScriptFinder(FolderManager.FolderManager folderManager)
    {
        _folderManager = folderManager;
    }

    public string? GetBuildScriptToCopy()
    {
        var localScriptsFolder = _folderManager.GetLocalScriptsFolder();

        var machineLocalScriptsFolder = Path.Join(localScriptsFolder, Environment.MachineName);

        var buildScriptToCopy = Path.Join(machineLocalScriptsFolder, Environment.UserName, buildScript);

        if (File.Exists(buildScriptToCopy))
        {
            return buildScriptToCopy;
        }

        buildScriptToCopy = Path.Join(machineLocalScriptsFolder, buildScript);

        if (File.Exists(buildScriptToCopy))
        {
            return buildScriptToCopy;
        }

        return null;
    }

    public string GetBuildScriptDestination()
    {
        var autogenFolder = _folderManager.GetAutogenFolder();

        return Path.Join(autogenFolder, buildScript);
    }
}
