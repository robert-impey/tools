namespace Make2CMakeLists;

internal static class MakefileParser
{
    public static Makefile? Parse(FileInfo makefileFileInfo)
    {
        int? standard = null;
        string? src = null;
        string? target = null;

        foreach (var line in File.ReadLines(makefileFileInfo.FullName))
        {
            if (standard is null)
            {
                standard = ParseCxxFlagsLine(line);
            }

            if (src is null)
            {
                src = ParseVarLine("SRC", line);
            }

            if (target is null)
            {
                target = ParseVarLine("TARGET", line);
            }
        }

        if (standard is null || src is null || target is null)
            return null;

        return new Makefile
        {
            CppStandard = standard.Value,
            Src = src,
            Target = target
        };
    }

    public static int? ParseCxxFlagsLine(string line)
    {
        var cxxRight = ParseVarLine("CXXFLAGS", line);

        if (cxxRight is null)
            return null;

        var flags = cxxRight.Split(' ');

        if (flags.Length == 0)
            return null;

        const string stdLeft = "-std=c++";
        foreach (var flag in flags)
        {
            if (!flag.StartsWith(stdLeft)) continue;
            var std = flag[stdLeft.Length..];

            if (int.TryParse(std, out var value))
                return value;
        }

        return 17;
    }

    public static string? ParseVarLine(string varName, string line)
    {
        if (!line.StartsWith(varName))
            return null;

        var lineAfterName = line[varName.Length..].TrimStart();

        if (lineAfterName.Length == 0 || lineAfterName[0] != '=')
            return null;

        var lineAfterAssignment = lineAfterName[1..].TrimStart();

        return lineAfterAssignment.Length > 0 ? lineAfterAssignment : null;
    }
}