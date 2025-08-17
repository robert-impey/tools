using Microsoft.Extensions.Logging;

namespace FolderManager;

public class LinuxFolderManager(ILogger<FolderManager> logger) : UnixFolderManager(logger)
{
    protected override string GetLocationsFile() => GetLocationsFile("linux");
}
