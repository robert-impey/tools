using Microsoft.Extensions.Logging;

namespace FolderManager;

public abstract class UnixFolderManager(ILogger<FolderManager> logger) : FolderManager(logger)
{
    public override string GetLocalScriptsFolder()
    {
        var localScriptsPathParts = new List<string>();

        var home = Environment.GetEnvironmentVariable("HOME");

        if (home is not null)
        {
            localScriptsPathParts.Add(home);
            localScriptsPathParts.Add("local-scripts");
        }

        if (localScriptsPathParts.Count == 0)
        {
            throw new ApplicationException("Unable to find the local scripts folder!");
        }

        var localScriptsFolder = Path.Combine(localScriptsPathParts.ToArray());
        var symlink = new FileInfo(localScriptsFolder);
        var target = symlink.ResolveLinkTarget(true); // true = follow to final target

        return target is null ? localScriptsFolder : target.FullName;
    }
}