using Shouldly;

namespace ResetPerms.Test;

public class ScriptsFinderTests
{
    [Fact]
    public void ShouldFindUnixScript()
    {
        ScriptsFinder.FileHasShebang("Files/unix-script.sh").ShouldBeTrue();
    }

    [Theory]
    [InlineData("Files/Empty.txt")]
    [InlineData("Files/NotScript.txt")]
    [InlineData("Files/WindowsOnly.ps1")]
    public void ShouldIgnoreNonScripts(string fileName)
    {
        ScriptsFinder.FileHasShebang(fileName).ShouldBeFalse();
    }
}
