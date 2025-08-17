using Microsoft.Extensions.Logging;

namespace FolderManager;

public class MacOsFolderManager(ILogger<FolderManager> logger) : UnixFolderManager(logger)
{
    protected override string GetLocationsFile() => GetLocationsFile("darwin");
}