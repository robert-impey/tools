var folderManager = FolderManager.FolderManager.GetFolderManager();

var managedFolders = await folderManager.GetManagedFolders();
var locations = managedFolders.Keys.ToArray();
Array.Sort(locations);

var autogenFolder = folderManager.GetAutogenFolder();

var managedFoldersFile = Path.Combine(autogenFolder, "managed-folders.txt");

using var outputFile = new StreamWriter(managedFoldersFile, false);

await outputFile.WriteLineAsync("# AUTOGEN'D - DO NOT EDIT!\n"
                                + $"# Written {DateTime.UtcNow:u}\n");

var havePrinted = false;
foreach (var location in locations)
{
    var foldersInLocation = managedFolders[location].ToArray();

    if (foldersInLocation.Length == 0)
    {
        continue;
    }

    if (havePrinted)
    {
        await outputFile.WriteLineAsync();
    }

    Array.Sort(foldersInLocation);

    foreach (var folder in foldersInLocation)
    {
        await outputFile.WriteLineAsync(folder);
    }

    havePrinted = true;
}
