using FolderManager;
using Microsoft.Extensions.Logging;
using System.Globalization;

namespace GenerateWindowsConfigSynchScripts;

internal class WindowsConfigScriptsGenerator
{
    private readonly ILogger<WindowsConfigScriptsGenerator> _logger;

    private readonly string _logsDirectory, _autogen, _files, _source, _destination;

    public WindowsConfigScriptsGenerator(
        ILogger<WindowsConfigScriptsGenerator> logger,
        string logsDirectory,
        string autogen,
        string files,
        string source,
        string destination
        )
    {
        ArgumentNullException.ThrowIfNull(logger);
        ArgumentNullException.ThrowIfNull(logsDirectory);
        ArgumentNullException.ThrowIfNull(autogen);
        ArgumentNullException.ThrowIfNull(files);
        ArgumentNullException.ThrowIfNull(source);
        ArgumentNullException.ThrowIfNull(destination);

        _logger = logger;
        _logsDirectory = logsDirectory;
        _autogen = autogen;
        _files = files;
        _source = source;
        _destination = destination;

        _logger.LogInformation("Creating WindowsConfigScriptsGenerator");

        _logger.LogInformation($"LogsDirectory: {_logsDirectory}");
        _logger.LogInformation($"Autogen: {_autogen}");
        _logger.LogInformation($"Files: {_files}");
        _logger.LogInformation($"Source: {_source}");
        _logger.LogInformation($"Destination: {_destination}");
    }

    public async Task Generate()
    {
        _logger.LogInformation("Generating scripts...");
        
        var outputScriptPath = Path.Combine(_autogen, "config-Windows.ps1");

        if (File.Exists(outputScriptPath))
        {
            _logger.LogInformation($"Deleting existing script at {outputScriptPath}");
            File.Delete(outputScriptPath);
        }

        using var outputScriptWriter = new StreamWriter(outputScriptPath);
        await outputScriptWriter.WriteLineAsync("# AUTOGEN'D - DO NOT EDIT!");

        // Get current local time with offset
        var zonedDateTime = DateTimeOffset.Now;

        // Format using RFC1123
        var formattedDateTime = zonedDateTime.ToString("R", CultureInfo.InvariantCulture); // "R" = RFC1123

        await outputScriptWriter.WriteLineAsync($"# Written {formattedDateTime}");
        await outputScriptWriter.WriteLineAsync();

        await Task.CompletedTask;
    }
}
