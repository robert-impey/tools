using FluentAssertions;
using Microsoft.Extensions.Logging.Abstractions;

namespace FolderManager.Test;

public class FolderManagerTest
{
    [Fact]
    public void CanGetFolderManager()
    {
        var folderManager = FolderManager.GetFolderManager(NullLogger<FolderManager>.Instance);

        folderManager.Should().NotBeNull();
    }
}