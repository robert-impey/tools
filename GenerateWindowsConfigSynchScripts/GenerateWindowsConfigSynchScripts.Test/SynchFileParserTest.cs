using Shouldly;

namespace GenerateWindowsConfigSynchScripts.Test;

public class SynchFileParserTest
{
    [Fact]
    public async Task ShouldReturnCorrectSourceAndDestinationPaths()
    {
        // Arrange
        var filePath = @"Files\Fruit.txt";
        
        // Act
        var synchFile = await SynchFileParser.ParseFile(filePath);

        // Assert
        synchFile.ShouldNotBeNull();

        synchFile.Source.ShouldBe(@"C:\");
        synchFile.Destination.ShouldBe(@"D:\");

        synchFile.Files.ShouldNotBeNull();
        synchFile.Files.Count().ShouldBe(3);

        synchFile.Files.ShouldContain("apples.txt");
        synchFile.Files.ShouldContain("bananas.docx");
        synchFile.Files.ShouldContain("cherries.pdf");
    }
}
