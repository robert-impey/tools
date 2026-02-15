using Microsoft.Extensions.Logging.Abstractions;
using Shouldly;

namespace LogsDeleter.Test;

public class DeleteFromTests : IDisposable
{
    private readonly string _tempDir;
    private readonly Sweeping.LogsDeleter _logsDeleter;

    public DeleteFromTests()
    {
        _tempDir = Path.Combine(Path.GetTempPath(), "DeleteFromTests_" + Guid.NewGuid().ToString());
        Directory.CreateDirectory(_tempDir);
        _logsDeleter = new Sweeping.LogsDeleter(NullLogger<Sweeping.LogsDeleter>.Instance);
    }

    public void Dispose()
    {
        if (Directory.Exists(_tempDir))
            Directory.Delete(_tempDir, true);
    }

    private string CreateFile(string name, int size, DateTime modTime)
    {
        string path = Path.Combine(_tempDir, name);
        File.WriteAllBytes(path, new byte[size]);
        File.SetLastWriteTime(path, modTime);
        return path;
    }

    [Fact]
    public void TestDeleteFrom_DeletesOldFiles()
    {
        // Arrange
        int days = 7;
        var oldFile = CreateFile("old.log", 10, DateTime.Now.AddDays(-(days + 1)));
        var recentFile = CreateFile("recent.log", 10, DateTime.Now.AddHours(-1));

        // Act
        _logsDeleter.DeleteFrom(_tempDir, days, false, false);

        // Assert
        File.Exists(oldFile).ShouldBeFalse();
        File.Exists(recentFile).ShouldBeTrue();
    }

    [Fact]
    public void TestDeleteFrom_DeletesEmptyWhenFlag()
    {
        // Arrange
        int days = 30;
        var emptyRecent = CreateFile("empty.txt", 0, DateTime.Now.AddMinutes(-30));
        var nonEmptyRecent = CreateFile("data.txt", 5, DateTime.Now.AddMinutes(-30));

        // Act
        _logsDeleter.DeleteFrom(_tempDir, days, true, false);

        // Assert
        File.Exists(emptyRecent).ShouldBeFalse();
        File.Exists(nonEmptyRecent).ShouldBeTrue();
    }

    [Fact]
    public void TestDeleteFrom_DoesNotDeleteRecentWhenFlagFalse()
    {
        // Arrange
        int days = 10;
        var emptyRecent = CreateFile("empty.txt", 0, DateTime.Now.AddMinutes(-10));
        var nonEmptyRecent = CreateFile("data.txt", 5, DateTime.Now.AddMinutes(-10));

        // Act
        _logsDeleter.DeleteFrom(_tempDir, days, false, false);

        // Assert
        File.Exists(emptyRecent).ShouldBeTrue();
        File.Exists(nonEmptyRecent).ShouldBeTrue();
    }

    [Fact]
    public void TestDeleteFrom_WithNegativeDays_FailsValidation()
    {
        // Arrange
        int days = -1;
        var recentFile = CreateFile("recent.log", 10, DateTime.Now.AddMinutes(-1));

        // Act & Assert
        // Shouldly's way of handling exceptions
        Should.Throw<ArgumentException>(() =>
            _logsDeleter.DeleteFrom(_tempDir, days, false, false)
        );

        File.Exists(recentFile).ShouldBeTrue();
    }
}
