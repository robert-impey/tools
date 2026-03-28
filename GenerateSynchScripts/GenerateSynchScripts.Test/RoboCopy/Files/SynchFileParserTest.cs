using GenerateSynchScripts.RoboCopy.Files;
using Shouldly;

namespace GenerateSynchScripts.Test.RoboCopy.Files;

public class SynchFileParserTest
{
    [Fact]
    public async Task ParseFile_ShouldReturnCorrectSourceAndDestinationPaths()
    {
        // Arrange
        var filePath = "TestData/RoboCopy/Files/Fruit.txt";

        // Act
        var synchFile = await SynchFileParser.ParseFile(filePath);

        // Assert
        synchFile.ShouldNotBeNull();

        synchFile.Id.ShouldBe("Fruit");

        synchFile.Source.ShouldBe(@"C:\");
        synchFile.Destination.ShouldBe(@"D:\");

        synchFile.Files.ShouldNotBeNull();
        synchFile.Files.Count().ShouldBe(3);

        synchFile.Files.ShouldContain("apples.txt");
        synchFile.Files.ShouldContain("bananas.docx");
        synchFile.Files.ShouldContain("cherries.pdf");
    }

    [Theory]
    [InlineData("TestData/RoboCopy/Files/NoFiles.txt")]
    [InlineData("TestData/RoboCopy/Files/NoBlankLine.txt")]
    [InlineData("TestData/Empty.txt")]
    public async Task ParseFile_ShouldThrowIfNoFiles(string filePath)
    {
        await Should.ThrowAsync<InvalidOperationException>(async () =>
        {
            await SynchFileParser.ParseFile(filePath);
        });
    }
}
