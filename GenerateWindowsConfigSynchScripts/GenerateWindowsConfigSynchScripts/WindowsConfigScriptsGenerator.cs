using System.Text;
using Microsoft.Extensions.Logging;

namespace GenerateWindowsConfigSynchScripts;

internal class WindowsConfigScriptsGenerator
{
    private readonly ILogger<WindowsConfigScriptsGenerator> _logger;

    private readonly string _autogen, _script, _files, _source, _destination;

    public WindowsConfigScriptsGenerator(
        ILogger<WindowsConfigScriptsGenerator> logger,
        string logsDirectory,
        string autogen,
        string script,
        string files,
        string source,
        string destination
        )
    {
        ArgumentNullException.ThrowIfNull(logger);
        ArgumentNullException.ThrowIfNull(logsDirectory);
        ArgumentNullException.ThrowIfNull(autogen);
        ArgumentNullException.ThrowIfNull(script);
        ArgumentNullException.ThrowIfNull(files);
        ArgumentNullException.ThrowIfNull(source);
        ArgumentNullException.ThrowIfNull(destination);

        _logger = logger;
        _autogen = autogen;
        _script = script;
        _files = files;
        _source = source;
        _destination = destination;

        _logger.LogInformation("Creating WindowsConfigScriptsGenerator");

        _logger.LogInformation($"LogsDirectory: {logsDirectory}");
        _logger.LogInformation($"Autogen: {_autogen}");
        _logger.LogInformation($"Script: {_script}");
        _logger.LogInformation($"Files: {_files}");
        _logger.LogInformation($"Source: {_source}");
        _logger.LogInformation($"Destination: {_destination}");
    }

    public async Task Generate()
    {
        _logger.LogInformation("Generating scripts...");
        
        var outputScriptPath = Path.Combine(_autogen, $"{_script}.ps1");

        if (File.Exists(outputScriptPath))
        {
            _logger.LogInformation($"Deleting existing script at {outputScriptPath}");
            File.Delete(outputScriptPath);
        }

        var sb = new StringBuilder();
        sb.Append("# AUTOGEN'D - DO NOT EDIT!\n");

        sb.Append($"# Written {DateTimeOffset.Now:R}\n\n");

        var first = true;
        foreach (var file in await File.ReadAllLinesAsync(_files))
        {
            if (first)
            {
                first = false;
            }
            else
            {
                sb.Append('\n');
            }
            
            if (string.IsNullOrWhiteSpace(file) || file.StartsWith('#'))
            {
                continue;
            }

            sb.Append($"ROBOCOPY \"{_source}\" \"{_destination}\" /xo {file}\n");
            sb.Append($"ROBOCOPY \"{_destination}\" \"{_source}\" /xo {file}\n");
        }

        await using var outputScriptWriter = new StreamWriter(outputScriptPath);
        await outputScriptWriter.WriteAsync(sb);
    }
}
