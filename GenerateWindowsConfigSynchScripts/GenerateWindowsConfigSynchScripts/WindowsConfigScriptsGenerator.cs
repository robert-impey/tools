using System.Text;
using Microsoft.Extensions.Logging;

namespace GenerateWindowsConfigSynchScripts;

internal class WindowsConfigScriptsGenerator
{
    private readonly ILogger<WindowsConfigScriptsGenerator> _logger;

    public WindowsConfigScriptsGenerator(
        ILogger<WindowsConfigScriptsGenerator> logger
        )
    {
        ArgumentNullException.ThrowIfNull(logger);


        _logger = logger;


        _logger.LogInformation("Creating WindowsConfigScriptsGenerator");
    }

    public async Task Generate(
        string autogen,
        string script,
        string source,
        string destination,
        IEnumerable<string> files
        )
    {
        ArgumentNullException.ThrowIfNull(autogen);
        ArgumentNullException.ThrowIfNull(script);
        ArgumentNullException.ThrowIfNull(source);
        ArgumentNullException.ThrowIfNull(destination);
        ArgumentNullException.ThrowIfNull(files);

        _logger.LogInformation("Generating scripts...");
        _logger.LogInformation($"Autogen: {autogen}");
        _logger.LogInformation($"Script: {script}");
        _logger.LogInformation($"Source: {source}");
        _logger.LogInformation($"Destination: {destination}");
        _logger.LogInformation($"Files: {string.Join(", ", files)}");

        var outputScriptPath = Path.Combine(autogen, $"{script}.ps1");

        if (File.Exists(outputScriptPath))
        {
            _logger.LogInformation($"Deleting existing script at {outputScriptPath}");
            File.Delete(outputScriptPath);
        }

        var sb = new StringBuilder();
        sb.Append("# AUTOGEN'D - DO NOT EDIT!\n");

        sb.Append($"# Written {DateTimeOffset.Now:R}\n\n");

        var first = true;
        foreach (var file in files)
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

            sb.Append($"ROBOCOPY \"{source}\" \"{destination}\" /xo {file}\n");
            sb.Append($"ROBOCOPY \"{destination}\" \"{source}\" /xo {file}\n");
        }

        await using var outputScriptWriter = new StreamWriter(outputScriptPath);
        await outputScriptWriter.WriteAsync(sb);
    }
}
