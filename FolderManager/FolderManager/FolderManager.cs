using NLog;

namespace FolderManager;

public abstract class FolderManager
{
    protected readonly ILogger _logger;

    protected FolderManager(ILogger logger)
    {
        _logger = logger;
    }

    public static FolderManager GetFolderManager(ILogger? logger = null)
    {
        if (logger is null)
        {
            var config = new NLog.Config.LoggingConfiguration();

            var logconsole = new NLog.Targets.ConsoleTarget("logconsole");

            config.AddRule(LogLevel.Info, LogLevel.Fatal, logconsole);

            LogManager.Configuration = config;

            logger = LogManager.GetCurrentClassLogger();
        }

        if (OperatingSystem.IsWindows())
        {
            return new WindowsFolderManager(logger);
        }

        if (OperatingSystem.IsLinux())
        {
            return new LinuxFolderManager(logger);
        }

        throw new ApplicationException("No folder manager for your operating system!");
    }

    protected abstract string GetLocationsFile();
    public abstract string GetLocalScriptsFolder();
    protected static string HomeFolder
    {
        get
        {
            if (OperatingSystem.IsWindows())
            {
                return Environment.GetEnvironmentVariable("USERPROFILE") ?? throw new ApplicationException("USERPROFILE environment variable not set!");
            }

            if (OperatingSystem.IsLinux())
            {
                return Environment.GetEnvironmentVariable("HOME") ?? throw new ApplicationException("HOME environment variable not set!");
            }

            throw new ApplicationException("No home folder for your operating system!");
        }
    }

    public string GetCommonLocalScriptsFolder() => Path.Join(GetLocalScriptsFolder(), "_Common");

    public string GetAutogenFolder(bool ensureExists = true)
    {
        var autogenFolder = Path.Join(HomeFolder, "autogen");

        if (ensureExists && !Directory.Exists(autogenFolder))
        {
            Directory.CreateDirectory(autogenFolder);
        }

        return autogenFolder;
    }

    public string GetFoldersFile()
    {
        var foldersFile = Path.Join(GetCommonLocalScriptsFolder(), "folders.txt");
        
        _logger.Info($"Folders file: {foldersFile}");

        return foldersFile;
    }

    protected string GetLocationsFile(string operatingSystemPathPart)
    {
        var locationsPathParts = new List<string>
        {
            GetCommonLocalScriptsFolder(),
            operatingSystemPathPart,
            "locations.txt"
        };

        var locationsFile = Path.Join(locationsPathParts.ToArray());

        _logger.Info($"Locations file: {locationsFile}");

        return locationsFile;
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

    public static string LogsFolder => Path.Combine(HomeFolder, "logs");

    public abstract string PowerShellExe { get; }
}