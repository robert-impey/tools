using Microsoft.Extensions.Logging.Abstractions;
using Shouldly;
namespace GenerateWindowsConfigSynchScripts.Test;

public class ScriptGeneratorTest
{
    private readonly WindowsConfigScriptsGenerator _sut;

    public ScriptGeneratorTest()
    {
        _sut = new WindowsConfigScriptsGenerator(NullLogger<WindowsConfigScriptsGenerator>.Instance);
    }

    [Fact]
    public async Task Generate_ShouldThrow_WhenAutogenIsNull()
    {
#pragma warning disable CS8625 // Cannot convert null literal to non-nullable reference type
        await Should.ThrowAsync<ArgumentNullException>(() =>
            _sut.Generate("id", null, "script", "source", "destination", ["file1"]));
#pragma warning restore CS8625
    }

    [Fact]
    public async Task Generate_ShouldCreateScriptFile_WithExpectedContent()
    {
        // Arrange
        var tempDir = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString());
        Directory.CreateDirectory(tempDir);

        var autogen = tempDir;
        var scriptName = "myscript";
        var source = "C:\\Source";
        var destination = "C:\\Destination";
        var files = new[] { "file1.txt", "file2.txt" };

        var outputScriptPath = Path.Combine(autogen, $"{scriptName}.ps1");

        // Act
        await _sut.Generate("id123", autogen, scriptName, source, destination, files);

        // Assert
        File.Exists(outputScriptPath).ShouldBeTrue();

        var content = await File.ReadAllTextAsync(outputScriptPath);
        content.ShouldContain("# AUTOGEN'D - DO NOT EDIT!");
        content.ShouldContain("SynchSingleFile2Ways");
        content.ShouldContain("file1.txt");
        content.ShouldContain("file2.txt");
    }

    [Fact]
    public async Task Generate_ShouldOverwriteExistingScriptFile()
    {
        // Arrange
        var tempDir = Path.Combine(Path.GetTempPath(), Guid.NewGuid().ToString());
        Directory.CreateDirectory(tempDir);

        var autogen = tempDir;
        var scriptName = "myscript";
        var source = "C:\\Source";
        var destination = "C:\\Destination";
        var files = new[] { "file1.txt" };

        var outputScriptPath = Path.Combine(autogen, $"{scriptName}.ps1");

        // Create an existing file
        await File.WriteAllTextAsync(outputScriptPath, "Old content");

        // Act
        await _sut.Generate("id123", autogen, scriptName, source, destination, files);

        // Assert
        var content = await File.ReadAllTextAsync(outputScriptPath);
        content.ShouldNotContain("Old content");
        content.ShouldContain("SynchSingleFile2Ways");
    }
}
