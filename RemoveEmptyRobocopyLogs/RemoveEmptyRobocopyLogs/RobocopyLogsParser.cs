using System.Text.RegularExpressions;

namespace RemoveEmptyRobocopyLogs;

public static partial class RobocopyLogsParser
{
    private readonly static Regex filesLineRegex = FilesLineRegex();

    public static bool IsFilesCopiedLine(string line)
    {
        var match = filesLineRegex.Match(line);

        if (match.Success)
        {
            var copiedCountString = match.Groups[1].Value;

            if (int.TryParse(copiedCountString, out var copiedCount))
            {
                return copiedCount > 0;
            }
        }

        return false;
    }

    public async static Task<bool> FileHasCopies(string fileName, CancellationToken cancellationToken)
    {
        var lines = await File.ReadAllLinesAsync(fileName, cancellationToken);

        // iterate backwards
        for (var i = lines.Length - 1; i >= 0; i--)
        {
            cancellationToken.ThrowIfCancellationRequested();

            if (IsFilesCopiedLine(lines[i]))
            {
                return true; // found a match near the end
            }
        }

        return false; // no match found
    }

    [GeneratedRegex(@"\s*Files\s*:\s+\d+\s+(\d+)\s+")]
    private static partial Regex FilesLineRegex();
}
