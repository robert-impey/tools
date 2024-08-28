namespace GenerateBuildScripts;

internal class BuildScriptFinder(FolderManager.FolderManager folderManager)
{
    private const string BuildScript = "Build.ps1";

    public string? GetBuildScriptToCopy()
    {
        var localScriptsFolder = folderManager.GetLocalScriptsFolder();

        var machineLocalScriptsFolder = Path.Join(localScriptsFolder, Environment.MachineName);

        var buildScriptToCopy = Path.Join(machineLocalScriptsFolder, Environment.UserName, BuildScript);

        if (File.Exists(buildScriptToCopy))
        {
            return buildScriptToCopy;
        }

        buildScriptToCopy = Path.Join(machineLocalScriptsFolder, BuildScript);

        if (File.Exists(buildScriptToCopy))
        {
            return buildScriptToCopy;
        }

        return null;
    }

    public string GetBuildScriptDestination()
    {
        var autogenFolder = folderManager.GetAutogenFolder();

        return Path.Join(autogenFolder, BuildScript);
    }
}
