using System.Text;
using NSubstitute;
using Shouldly;
using TidyFolder.Search;

namespace TidyFolder.Test;

public class SearchLogicTest
{
    [Fact]
    public void BuildDirsAndFiles_FindsAllFiles()
    {
        var tempDir = Directory.CreateTempSubdirectory();
        try
        {
            var dir1 = Path.Combine(tempDir.FullName, "dir1");
            var dir2 = Path.Combine(tempDir.FullName, "dir2");
            Directory.CreateDirectory(dir1);
            Directory.CreateDirectory(dir2);

            CreateFile(Path.Combine(dir1, "file1.txt"), "content1");
            CreateFile(Path.Combine(dir1, "file2.txt"), "content2");
            CreateFile(Path.Combine(dir2, "file3.txt"), "content3");
            CreateFile(Path.Combine(tempDir.FullName, "file4.txt"), "content4");

            var entries = SearchLogic.BuildDirsAndFiles(tempDir.FullName);

            entries.Count.ShouldBe(4);

            var names = entries.Select(e => e.Name).ToHashSet();
            names.ShouldContain("file1.txt");
            names.ShouldContain("file2.txt");
            names.ShouldContain("file3.txt");
            names.ShouldContain("file4.txt");
        }
        finally
        {
            tempDir.Delete(true);
        }
    }

    private static void CreateFile(string path, string content) =>
        File.WriteAllText(path, content);

    [Fact]
    public void FindMatchingStems_FindsExpectedPair()
    {
        var entries = new List<DirEntry>
        {
            new() { Path = "/tmp/test0/a.txt", Name = "a.txt" },
            new() { Path = "/tmp/test0/a1.txt", Name = "a1.txt" },
            new() { Path = "/tmp/test0/b.txt", Name = "b.txt" },
            new() { Path = "/tmp/test0/c.txt", Name = "c.txt" },
            new() { Path = "/tmp/test0/a.md", Name = "a.md" },
            new() { Path = "/tmp/test1/c1.txt", Name = "c1.txt" },
            new() { Path = "/tmp/test0/", Name = string.Empty },
            new() { Path = "/tmp/test1/", Name = string.Empty },
        };

        var matches = SearchLogic.FindMatchingStems(entries);

        matches.ShouldHaveSingleItem();

        var (first, second) = matches[0];

        Path.GetExtension(first.Name).ShouldBe(".txt");
        Path.GetFileNameWithoutExtension(second.Name)
            .StartsWith(Path.GetFileNameWithoutExtension(first.Name))
            .ShouldBeTrue();
    }

    [Fact]
    public void FindMatchingStems_DoesNotMatchFilesInCaseDifferentFolders()
    {
        var root = Path.Combine(Path.GetTempPath(), "TidyFolderTest");
        var entries = new List<DirEntry>
        {
            new() { Path = Path.Combine(root, "foo", "readme.md"), Name = "readme.md" },
            new() { Path = Path.Combine(root, "FOO", "readme-v2.md"), Name = "readme-v2.md" },
        };

        var matches = SearchLogic.FindMatchingStems(entries);

        matches.ShouldBeEmpty();
    }

    [Fact]
    public void FindMatchingStems_MatchesStemsRegardlessOfCaseInSameFolder()
    {
        var directory = Path.Combine(Path.GetTempPath(), "TidyFolderTest", "foo");
        var entries = new List<DirEntry>
        {
            new() { Path = Path.Combine(directory, "readme.md"), Name = "readme.md" },
            new() { Path = Path.Combine(directory, "ReadMe-v2.md"), Name = "ReadMe-v2.md" },
        };

        var matches = SearchLogic.FindMatchingStems(entries);

        matches.ShouldHaveSingleItem();
        matches[0].Item1.Name.ShouldBe("readme.md");
        matches[0].Item2.Name.ShouldBe("ReadMe-v2.md");
    }

    [Fact]
    public void FindMatchingStems_MatchesExtensionsRegardlessOfCase()
    {
        var directory = Path.Combine(Path.GetTempPath(), "TidyFolderTest", "foo");
        var entries = new List<DirEntry>
        {
            new() { Path = Path.Combine(directory, "report.md"), Name = "report.md" },
            new() { Path = Path.Combine(directory, "report-final.MD"), Name = "report-final.MD" },
        };

        var matches = SearchLogic.FindMatchingStems(entries);

        matches.ShouldHaveSingleItem();
        matches[0].Item1.Name.ShouldBe("report.md");
        matches[0].Item2.Name.ShouldBe("report-final.MD");
    }

    [Fact]
    public void FindMatchingStems_MatchesExtensionlessFiles()
    {
        var directory = Path.Combine(Path.GetTempPath(), "TidyFolderTest", "foo");
        var entries = new List<DirEntry>
        {
            new() { Path = Path.Combine(directory, "README"), Name = "README" },
            new() { Path = Path.Combine(directory, "README-v2"), Name = "README-v2" },
        };

        var matches = SearchLogic.FindMatchingStems(entries);

        matches.ShouldHaveSingleItem();
        matches[0].Item1.Name.ShouldBe("README");
        matches[0].Item2.Name.ShouldBe("README-v2");
    }

    [Fact]
    public void SearchDirectory_NoLogsDir_DoesNotError()
    {
        var tempDir = Directory.CreateTempSubdirectory();
        try
        {
            var dir1 = Path.Combine(tempDir.FullName, "dir1");
            Directory.CreateDirectory(dir1);

            CreateFile(Path.Combine(dir1, "file1.txt"), "1");
            CreateFile(Path.Combine(dir1, "file2.txt"), "2");

            var ex = Record.Exception(() =>
                SearchLogic.SearchDirectory(tempDir.FullName, ""));

            ex.ShouldBeNull();
        }
        finally
        {
            tempDir.Delete(true);
        }
    }

    [Fact]
    public void SearchDirectory_WritesToStdout_WhenNoLogsDir()
    {
        var searchDir = Directory.CreateTempSubdirectory();
        try
        {
            CreateFile(Path.Combine(searchDir.FullName, "a.txt"), "x");
            CreateFile(Path.Combine(searchDir.FullName, "a1.txt"), "y");

            var sw = new StringWriter();
            var oldOut = Console.Out;
            Console.SetOut(sw);

            try
            {
                SearchLogic.SearchDirectory(searchDir.FullName, "");
            }
            finally
            {
                Console.SetOut(oldOut);
            }

            var output = sw.ToString();
            output.ShouldContain("--- Results for:");
        }
        finally
        {
            searchDir.Delete(true);
        }
    }

    [Fact]
    public void SearchDirectory_WritesResultsOnlyToLog_WhenLogsDirSet()
    {
        var root = Directory.CreateTempSubdirectory();
        try
        {
            var logsDir = Path.Combine(root.FullName, "logs");
            Directory.CreateDirectory(logsDir);

            CreateFile(Path.Combine(root.FullName, "a.txt"), "x");
            CreateFile(Path.Combine(root.FullName, "a1.txt"), "y");

            var sw = new StringWriter();
            var oldOut = Console.Out;
            Console.SetOut(sw);

            try
            {
                SearchLogic.SearchDirectory(root.FullName, logsDir);
            }
            finally
            {
                Console.SetOut(oldOut);
            }

            var stdout = sw.ToString();
            stdout.ShouldNotContain("--- Results for:");

            var log = Directory.EnumerateFiles(logsDir, "*.log").Single();
            var content = File.ReadAllText(log);
            content.ShouldContain("--- Results for:");
        }
        finally
        {
            root.Delete(true);
        }
    }

    [Fact]
    public void SearchDirectory_WritesExceptionToErrorLog_WhenConsoleOutThrows()
    {
        var root = Directory.CreateTempSubdirectory();
        try
        {
            var logsDir = Path.Combine(root.FullName, "logs");
            Directory.CreateDirectory(logsDir);

            CreateFile(Path.Combine(root.FullName, "a.txt"), "x");
            CreateFile(Path.Combine(root.FullName, "a1.txt"), "y");

            var oldOut = Console.Out;

            var writer = Substitute.For<TextWriter>();
            writer.Encoding.Returns(Encoding.UTF8);
            writer.WhenForAnyArgs(x => x.Write(Arg.Any<char>())).Do(_ => throw new NotSupportedException("Test writer"));
            writer.WhenForAnyArgs(x => x.Write(Arg.Any<string>())).Do(_ => throw new NotSupportedException("Test writer"));
            writer.WhenForAnyArgs(x => x.WriteLine(Arg.Any<string>())).Do(_ => throw new NotSupportedException("Test writer"));
            writer.WhenForAnyArgs(x => x.WriteLine()).Do(_ => throw new NotSupportedException("Test writer"));
            writer.WhenForAnyArgs(x => x.Flush()).Do(_ => throw new NotSupportedException("Test writer"));

            Console.SetOut(writer);

            try
            {
                var ex = Record.Exception(() =>
                    SearchLogic.SearchDirectory(root.FullName, logsDir));

                ex.ShouldNotBeNull();
                ex.ShouldBeOfType<NotSupportedException>();
            }
            finally
            {
                Console.SetOut(oldOut);
            }

            var errLog = Directory.EnumerateFiles(logsDir, "*.err").Single();
            var content = File.ReadAllText(errLog);
            content.ShouldContain("ERROR processing");
            content.ShouldContain(root.FullName);
        }
        finally
        {
            root.Delete(true);
        }
    }


}
