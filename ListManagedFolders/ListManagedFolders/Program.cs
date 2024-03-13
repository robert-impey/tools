var folderManager = FolderManager.FolderManager.GetFolderManager();

var managedFolders = await folderManager.GetManagedFolders();
var locations = managedFolders.Keys.ToArray();
Array.Sort(locations);

var havePrinted = false;
for (var i = 0; i < locations.Length; i++)
{
    var location = locations[i];

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
