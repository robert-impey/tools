namespace GenerateSynchScripts;

public static class SynchFileParser
{
    public async static Task<SynchFile> ParseFile(string filePath)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(filePath);

        var id = Path.GetFileNameWithoutExtension(filePath);

        var lines = await File.ReadAllLinesAsync(filePath);

        if (lines.Length < 4)
        {
            throw new InvalidOperationException(
                "The file must contain at least four lines: source path, destination path, a blank line, and at least one file.");
        }

        var source = lines[0].Trim();
        var destination = lines[1].Trim();
        var files = lines.Skip(3).Select(line => line.Trim()).Where(line => !string.IsNullOrWhiteSpace(line)).ToList();
        if (string.IsNullOrWhiteSpace(source))
        {
            throw new InvalidOperationException("The source path cannot be empty.");
        }

        if (string.IsNullOrWhiteSpace(destination))
        {
            throw new InvalidOperationException("The destination path cannot be empty.");
        }

        if (files.Count == 0)
        {
            throw new InvalidOperationException("At least one file must be specified.");
        }

        return new SynchFile(id, source, destination, files);
    }
}
