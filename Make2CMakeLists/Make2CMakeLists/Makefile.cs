namespace Make2CMakeLists;

internal record Makefile
{
    public int CppStandard { get; init; }
    public required string Src { get; init; }
    public required string Target { get; init; }
}