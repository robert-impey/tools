namespace FolderManager;

public abstract class FolderManager
{
    public static FolderManager GetFolderManager()
    {
        if (OperatingSystem.IsWindows())
        {
            return new WindowsFolderManager();
        }

        if (OperatingSystem.IsLinux())
        {
            return new LinuxFolderManager();
        }

        throw new ApplicationException("No folder manager for your operating system!");
    }

    protected abstract string GetLocationsFile();
    protected abstract string GetLocalScriptsFolder();
    protected abstract string GetHomeFolder();

    public string GetCommonLocalScriptsFolder() => Path.Join(GetLocalScriptsFolder(), "_Common");

    public string GetFoldersFile() => Path.Join(GetCommonLocalScriptsFolder(), "folders.txt");

    protected string GetLocationsFile(string operatingSystemPathPart)
    {
        var locationsPathParts = new List<string>
        {
            GetCommonLocalScriptsFolder(),
            operatingSystemPathPart,
            "locations.txt"
        };

        return Path.Join(locationsPathParts.ToArray());
    }

    public async Task<Dictionary<string, List<string>>> GetManagedFolders()
    {
        var foldersPath = GetFoldersFile();

        if (!Path.Exists(foldersPath))
        {
            throw new ApplicationException($"{foldersPath} does not exist!");
        }

        var locationsPath = GetLocationsFile();

        if (!Path.Exists(locationsPath))
        {
            throw new ApplicationException($"{locationsPath} does not exist!");
        }

        var folders = await File.ReadAllLinesAsync(foldersPath);

        if (folders is null)
        {
            throw new ApplicationException($"Unable to read {foldersPath}!");
        }

        var locations = await File.ReadAllLinesAsync(locationsPath);
        if (locations is null)
        {
            throw new ApplicationException($"Unable to read {locationsPath}!");
        }

        var managedFolders = new Dictionary<string, List<string>>();

        foreach (var location in locations)
        {
            var locationFolders = new List<string>();

            foreach (var folder in folders)
            {
                var managedFolderPath = Path.Combine(location, folder);

                if (Directory.Exists(managedFolderPath))
                {
                    locationFolders.Add(managedFolderPath);
                }
            }

            managedFolders[location] = locationFolders;
        }

        return managedFolders;
    }

    public string? GetBuildScriptToRun()
    {
        var localScriptsFolder = GetLocalScriptsFolder();

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

    public string GetLogsFolder() => Path.Combine(GetHomeFolder(), "logs");
}