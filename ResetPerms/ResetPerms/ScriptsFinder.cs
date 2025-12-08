namespace ResetPerms;

public static class ScriptsFinder
{
    public static bool FileHasShebang(string fileName)
    {
        using var reader = new StreamReader(fileName);

        var firstLine = reader.ReadLine();

        return firstLine != null && firstLine.StartsWith("#!");
    }

    public static IEnumerable<string> FindFilesWithShebang(string scriptsDir)
    {
        var matchingFiles = Directory.EnumerateFiles(
            scriptsDir,
            "*",
            SearchOption.AllDirectories
        );

        return matchingFiles.Where(FileHasShebang);
    }
}
