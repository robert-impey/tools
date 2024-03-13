var folderManager = FolderManager.FolderManager.GetFolderManager();

var managedFolders = await folderManager.GetManagedFolders();
var locations = managedFolders.Keys.ToArray();
Array.Sort(locations);

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
        Console.WriteLine();
    }

    Array.Sort(foldersInLocation);

    foreach (var folder in foldersInLocation)
    {
        Console.WriteLine(folder);
    }

    havePrinted = true;
}
