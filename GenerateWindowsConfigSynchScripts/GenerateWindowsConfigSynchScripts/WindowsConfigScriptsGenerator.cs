using Microsoft.Extensions.Logging;
using System.Text;

namespace GenerateWindowsConfigSynchScripts;

internal class WindowsConfigScriptsGenerator
{
    private readonly ILogger<WindowsConfigScriptsGenerator> _logger;

    private readonly string _autogen, _script, _source, _destination;
    private readonly IEnumerable<string> _files; 

    public WindowsConfigScriptsGenerator(
        ILogger<WindowsConfigScriptsGenerator> logger,
        string autogen,
        string script,
        string source,
        string destination,
        IEnumerable<string> files
        )
    {
        ArgumentNullException.ThrowIfNull(logger);
        ArgumentNullException.ThrowIfNull(autogen);
        ArgumentNullException.ThrowIfNull(script);
        ArgumentNullException.ThrowIfNull(source);
        ArgumentNullException.ThrowIfNull(destination);
        ArgumentNullException.ThrowIfNull(files);

        _logger = logger;
        _autogen = autogen;
        _script = script;
        _source = source;
        _destination = destination;
        _files = files;

        _logger.LogInformation("Creating WindowsConfigScriptsGenerator");

        _logger.LogInformation($"Autogen: {_autogen}");
        _logger.LogInformation($"Script: {_script}");
        _logger.LogInformation($"Source: {_source}");
        _logger.LogInformation($"Destination: {_destination}");
        _logger.LogInformation($"Files: {string.Join(", ", _files)}");
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
        foreach (var file in _files)
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
