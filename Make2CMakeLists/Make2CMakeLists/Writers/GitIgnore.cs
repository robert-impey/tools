namespace Make2CMakeLists.Writers;

internal static class GitIgnore
{
    public static bool Write(Makefile makefile, DirectoryInfo directoryInfo)
    {
        var gitIgnore = Path.Join(directoryInfo.FullName, ".gitignore");

        using var outFile = File.CreateText(gitIgnore);

        outFile.WriteLine(makefile.Target);

        return true;
    }
}