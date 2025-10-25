using System.Text.RegularExpressions;

namespace RemoveEmptyRobocopyLogs;

public static partial class RobocopyLogsParser
{
    private static readonly Regex filesLineRegex = FilesLineRegex();

    public static bool IsFilesCopiedLine(string line)
    {
        var match = filesLineRegex.Match(line);

        if (match.Success)
        {
            string copiedCountString = match.Groups[1].Value;

            if (int.TryParse(copiedCountString, out int copiedCount))
            {
                return copiedCount > 0;
            }
        }
        return false;
    }

    public static bool FileHasCopies(string fileName)
    {
        return File.ReadAllLines(fileName).Any(IsFilesCopiedLine);
    }

    [GeneratedRegex(@"\s*Files\s*:\s+\d+\s+(\d+)\s+")]
    private static partial Regex FilesLineRegex();
}