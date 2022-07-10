namespace Make2CMakeLists;

internal record Makefile
{
    public int CppStandard { get; init; }
    public string Src { get; init; }
    public string Target { get; init; }
}