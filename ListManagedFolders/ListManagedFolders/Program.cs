var folderManager = FolderManager.FolderManager.GetFolderManager();

foreach (var managedFolder in await folderManager.GetManagedFolders())
{
    Console.WriteLine(managedFolder);
}
