namespace TidyFolder.Search;

public static class SearchLogic
{
    public static List<DirEntry> BuildDirsAndFiles(string root)
    {
        var entries = new List<DirEntry>();

        try
        {
            foreach (var path in Directory.EnumerateFiles(
                         root, "*", SearchOption.AllDirectories))
            {
                try
                {
                    var info = new FileInfo(path);

                    entries.Add(new DirEntry
                    {
                        Path = path,
                        Name = info.Name,
                        IsDir = false,
                        ModTime = info.LastWriteTime,
                        Size = info.Length,
                    });
                }
                catch (Exception ex)
                {
                    Console.Error.WriteLine(
                        $"Warning: Could not get info for {path}: {ex.Message}");
                }
            }
        }
        catch (Exception ex)
        {
            throw new InvalidOperationException(
                $"Error walking directory {root}", ex);
        }

        return entries;
    }

    public static List<(DirEntry, DirEntry)>
        FindMatchingStems(List<DirEntry> dirsAndFiles)
    {
        var result = new List<(DirEntry, DirEntry)>();

        for (int i = 0; i < dirsAndFiles.Count; i++)
        {
            var file = dirsAndFiles[i];
            if (!SplitStemExt(file.Path, out var stem, out var ext))
                continue;

            for (int j = 0; j < dirsAndFiles.Count; j++)
            {
                if (i == j)
                    continue;

                var other = dirsAndFiles[j];
                if (!SplitStemExt(other.Path, out var otherStem, out var otherExt))
                    continue;

                // Normalize possible nulls returned by Path APIs before comparing.
                var aExt = ext ?? string.Empty;
                var bExt = otherExt ?? string.Empty;
                if (!StringComparer.OrdinalIgnoreCase.Equals(aExt, bExt))
                    continue;

                var dir1 = Path.GetDirectoryName(file.Path) ?? string.Empty;
                var dir2 = Path.GetDirectoryName(other.Path) ?? string.Empty;
                if (!dir1.Equals(dir2, StringComparison.OrdinalIgnoreCase))
                {
                    continue;
                }

                // Ensure stems are non-null before calling StartsWith to avoid CS8602.
                if (!string.IsNullOrEmpty(otherStem)
                    && !string.IsNullOrEmpty(stem)
                    && otherStem.StartsWith(stem, StringComparison.Ordinal)
                    && otherStem != stem)
                {
                    result.Add((file, other));
                }
            }
        }

        return result;
    }

    public static void SearchDirectory(string dir, string logsDir)
    {
        var dirsAndFiles = BuildDirsAndFiles(dir);
        var matches = FindMatchingStems(dirsAndFiles);

        if (matches.Count == 0)
            return;

        TextWriter output = Console.Out;
        StreamWriter? logFile = null;
        StreamWriter? errFile = null;

        try
        {
            if (!string.IsNullOrEmpty(logsDir))
            {
                var safe = SanitizeForFileName(dir);
                var timestamp = GetLogTime();

                var logPath = Path.Combine(
                    logsDir, $"{timestamp}-search-{safe}.log");
                var errPath = Path.Combine(
                    logsDir, $"{timestamp}-search-{safe}.err");

                try
                {
                    logFile = new StreamWriter(logPath, false, System.Text.Encoding.UTF8);
                    output = logFile;
                }
                catch (Exception ex)
                {
                    Console.Error.WriteLine(
                        $"Warning: Could not create log file {logPath}: {ex.Message}");
                }

                try
                {
                    errFile = new StreamWriter(errPath);
                }
                catch (Exception ex)
                {
                    Console.Error.WriteLine(
                        $"Warning: Could not create error log file {errPath}: {ex.Message}");
                }
            }

            PrintMatchingStems(output, dir, matches);

            if (logFile != null)
                Console.WriteLine($"OK: processed {dir}");
        }
        catch (Exception ex)
        {
            errFile?.WriteLine($"ERROR processing {dir}: {ex}");
            throw;
        }
        finally
        {
            logFile?.Dispose();
            errFile?.Dispose();
        }
    }

    private static void PrintMatchingStems(TextWriter output, string name,
        List<(DirEntry File, DirEntry Other)> matches)
    {
        if (matches.Count == 0)
            return;

        // This ensures the writer uses UTF-8 if it's a file stream
        // and hasn't been configured otherwise.
        output.WriteLine();
        output.WriteLine($"--- Results for: {name} ---");

        foreach (var (file, other) in matches)
        {
            var parent = Path.GetDirectoryName(file.Path) ?? string.Empty;
            output.WriteLine(parent);
            output.WriteLine($"\t{file.Name}");
            output.WriteLine($"\t{other.Name}");
        }

        output.Flush();
    }

    private static string GetLogTime() =>
        DateTime.Now.ToString("yyyy-MM-dd_HH.mm.ss");

    private static string SanitizeForFileName(string s) =>
        s.Replace("/", "_").Replace("\\", "_").Replace(":", "_");

    private static bool SplitStemExt(
        string path,
        out string stem,
        out string ext)
    {
        var file = Path.GetFileName(path);
        if (string.IsNullOrEmpty(file))
        {
            stem = ext = string.Empty;
            return false;
        }

        ext = Path.GetExtension(file);
        stem = Path.GetFileNameWithoutExtension(file);
        return !string.IsNullOrEmpty(stem);
    }
}
