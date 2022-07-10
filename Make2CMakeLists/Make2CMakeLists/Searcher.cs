namespace Make2CMakeLists;

internal static class Searcher
{
    public static IEnumerable<FileInfo> FindMakeFiles(DirectoryInfo directoryInfo)
    {
        return Directory.GetFiles(directoryInfo.FullName, "Makefile", SearchOption.AllDirectories)
            .Select(path => new FileInfo(path));
    }
}