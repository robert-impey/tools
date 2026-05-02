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

            Assert.Equal(4, entries.Count);

            var names = entries.Select(e => e.Name).ToHashSet();
            Assert.Contains("file1.txt", names);
            Assert.Contains("file2.txt", names);
            Assert.Contains("file3.txt", names);
            Assert.Contains("file4.txt", names);
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
        };

        var matches = SearchLogic.FindMatchingStems(entries);

        Assert.Single(matches);

        var (first, second) = matches[0];

        Assert.Equal(".txt", Path.GetExtension(first.Name));
        Assert.True(Path.GetFileNameWithoutExtension(second.Name)
            .StartsWith(Path.GetFileNameWithoutExtension(first.Name)));
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

            Assert.Null(ex);
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
            Assert.Contains("--- Results for:", output);
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
            Assert.DoesNotContain("--- Results for:", stdout);

            var log = Directory.EnumerateFiles(logsDir, "*.log").Single();
            var content = File.ReadAllText(log);
            Assert.Contains("--- Results for:", content);
        }
        finally
        {
            root.Delete(true);
        }
    }
}
