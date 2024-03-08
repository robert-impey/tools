﻿namespace FolderManager;

public class LinuxFolderManager : FolderManager
{
    protected override string GetLocalScriptsFolder()
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

        return Path.Join(localScriptsPathParts.ToArray());
    }

    protected override string GetLocationsFile() => GetLocationsFile("linux");
}
