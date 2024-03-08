using FluentAssertions;

namespace FolderManager.Test;

public class FolderManagerTest
{
    [Fact]
    public void CanGetFolderManager()
    {
        var folderManager = FolderManager.GetFolderManager();

        folderManager.Should().NotBeNull();
    }
}