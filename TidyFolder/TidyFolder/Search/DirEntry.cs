namespace TidyFolder.Search;

public sealed class DirEntry
{
    public string Path { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public bool IsDir { get; init; }
    public DateTime ModTime { get; init; }
    public long Size { get; init; }
}
