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

    public static async Task<bool> FileHasCopies(string fileName, CancellationToken cancellationToken)
    {
        var lines = await File.ReadAllLinesAsync(fileName, cancellationToken);

        // iterate backwards
        for (int i = lines.Length - 1; i >= 0; i--)
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
