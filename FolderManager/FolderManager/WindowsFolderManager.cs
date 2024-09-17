using NLog;

namespace FolderManager;

public class WindowsFolderManager(ILogger logger) : FolderManager(logger)
{
    public override string PowerShellExe => "PowerShell";

    public override string GetLocalScriptsFolder()
    {
        var localScriptsPathParts = new List<string>();

        var localScripts = Environment.GetEnvironmentVariable("LOCAL_SCRIPTS");

        if (localScripts is not null)
        {
            localScriptsPathParts.Add(localScripts);
        }

        if (localScriptsPathParts.Count == 0)
        {
            throw new ApplicationException("Unable to find the local scripts folder!");
        }

        return Path.Join(localScriptsPathParts.ToArray());
    }

    protected override string GetLocationsFile() => GetLocationsFile("Windows");
}
