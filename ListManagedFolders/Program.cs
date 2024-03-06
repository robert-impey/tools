internal class Program
{
    private static async Task Main(string[] args)
    {
        var commonLocalScriptsPathParts = new List<string>();

        if (OperatingSystem.IsWindows())
        {
            var localScripts = Environment.GetEnvironmentVariable("LOCAL_SCRIPTS");

            if (localScripts is not null)
            {
                commonLocalScriptsPathParts.Add(localScripts);
            }
        }

        if (OperatingSystem.IsLinux())
        {
            var home = Environment.GetEnvironmentVariable("HOME");

            if (home is not null)
            {
                commonLocalScriptsPathParts.Add(home);
                commonLocalScriptsPathParts.Add("local-scripts");
            }
        }

        if (commonLocalScriptsPathParts.Count == 0)
        {
            throw new ApplicationException("Unable to find the local scripts folder!");
        }

        commonLocalScriptsPathParts.Add("_Common");

        var foldersPathParts = new List<string>(commonLocalScriptsPathParts);

        foldersPathParts.Add("folders.txt");

        string foldersPath = Path.Join(foldersPathParts.ToArray());

        if (Path.Exists(foldersPath))
        {
            Console.WriteLine($"Found {foldersPath}");
        }
        else
        {
            throw new ApplicationException($"{foldersPath} does not exist!");
        }

        var locationsPathParts = new List<string>(commonLocalScriptsPathParts);

        if (OperatingSystem.IsWindows())
        {
            locationsPathParts.Add("Windows");
        }

        if (OperatingSystem.IsLinux())
        {
            locationsPathParts.Add("linux");
        }

        locationsPathParts.Add("locations.txt");

        var locationsPath = Path.Join(locationsPathParts.ToArray());

        if (Path.Exists(locationsPath))
        {
            Console.WriteLine($"Found {locationsPath}");
        }
        else
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

        var managedFolders = new List<string>();

        foreach (var location in locations)
        {
            foreach (var folder in folders)
            {
                var managedFolderPath = Path.Combine(location, folder);

                if (Directory.Exists(managedFolderPath))
                {
                    managedFolders.Add(managedFolderPath);
                }
            }
        }

        managedFolders.Sort();

        foreach (var managedFolder in managedFolders)
        {
            Console.WriteLine(managedFolder);
        }
    }
}