using Microsoft.Extensions.Logging;

namespace LogsDeleter;

using System;
using System.Collections.Generic;
using System.IO;

public class LogsDeleter
{
    private readonly ILogger<LogsDeleter> _logger;

    public LogsDeleter(ILogger<LogsDeleter> logger)
    {
        ArgumentNullException.ThrowIfNull(logger);

        _logger = logger;
    }

    public void DeleteFrom(string subPath, int days, bool deleteEmpty, bool verbose)
    {
        if (days <= 0)
        {
            throw new ArgumentException($"days must be positive, got: {days}", nameof(days));
        }

        // Go's AddDate(0, 0, -days) equivalent
        DateTime cutoff = DateTime.Now.AddDays(-days);

        if (verbose)
        {
            _logger.LogInformation($"Searching {subPath} for files older than {cutoff}");
        }

        if (!Directory.Exists(subPath))
        {
            throw new DirectoryNotFoundException($"The path {subPath} does not exist.");
        }

        // Using EnumerateFiles to get full paths, similar to filepath.Glob
        var allFiles = Directory.EnumerateFiles(subPath, "*");
        var filesToDelete = new List<FileInfo>();

        foreach (var filePath in allFiles)
        {
            var fileInfo = new FileInfo(filePath);

            // Logic: Older than cutoff OR (deleteEmpty is true AND file is 0 bytes)
            if (fileInfo.LastWriteTime < cutoff)
            {
                filesToDelete.Add(fileInfo);
            }
            else if (deleteEmpty && fileInfo.Length == 0)
            {
                filesToDelete.Add(fileInfo);
            }
        }

        if (verbose || filesToDelete.Count > 0)
        {
            _logger.LogInformation($"Found {filesToDelete.Count} files to delete in {subPath}");
        }

        foreach (var file in filesToDelete)
        {
            _logger.LogInformation($"Deleting {file.FullName}");

            try
            {
                // File.Delete is the equivalent of os.Remove
                // If you need to delete directories too, use Directory.Delete(path, true)
                file.Delete();
            }
            catch (Exception ex)
            {
                // Go returns the error; C# typically throws.
                // Wrapping in a custom exception to match the "return err" flow.
                throw new IOException($"Failed to delete {file.FullName}", ex);
            }
        }
    }
}
