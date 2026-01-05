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
        string id,
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

        var sourceClean = CleanFolderPathForLogName(source);
        var destinationClean = CleanFolderPathForLogName(destination);

        var outputScriptPath = Path.Combine(autogen, $"{script}.ps1");

        if (File.Exists(outputScriptPath))
        {
            _logger.LogInformation($"Deleting existing script at {outputScriptPath}");
            File.Delete(outputScriptPath);
        }

        var sb = new StringBuilder();
        sb.AppendLine("# AUTOGEN'D - DO NOT EDIT!");

        sb.AppendLine($"# Written {DateTimeOffset.Now:R}");
        sb.AppendLine();

        sb.AppendLine("param(");
        sb.AppendLine("    [Parameter (Mandatory = $False)]");
        sb.AppendLine("    [switch]$logged = $False");
        sb.AppendLine(")");
        sb.AppendLine();

        sb.AppendLine(@"Import-Module ""$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1""");
        sb.AppendLine();
        
        var first = true;
        foreach (var file in files)
        {
            if (first)
            {
                first = false;
            }
            else
            {
                sb.AppendLine();
            }

            if (string.IsNullOrWhiteSpace(file) || file.StartsWith('#'))
            {
                continue;
            }

            sb.AppendLine("SynchSingleFile2Ways `");
            sb.AppendLine($"    -id \"{id}\" `");
            sb.AppendLine($"    -file \"{file}\" `");
            sb.AppendLine($"    -sourceFolder \"{source}\" `");
            sb.AppendLine($"    -sourceLogName \"{sourceClean}\" `");
            sb.AppendLine($"    -destinationFolder \"{destination}\" `");
            sb.AppendLine($"    -destinationLogName \"{destinationClean}\" `");
            sb.AppendLine("    -logged $logged");
        }

        await using var outputScriptWriter = new StreamWriter(outputScriptPath);
        await outputScriptWriter.WriteAsync(sb);
    }

    private static string CleanFolderPathForLogName(string path)
    {
        if (string.IsNullOrWhiteSpace(path))
            throw new ArgumentException("Path cannot be null or empty.", nameof(path));

        // Normalise slashes
        var normalised = path.Replace('/', '\\');

        // Replace backslashes with underscores
        var replaced = normalised.Replace('\\', '_');

        // Remove invalid filename characters
        var invalid = Path.GetInvalidFileNameChars();
        var sb = new StringBuilder(replaced.Length);

        foreach (var ch in replaced)
        {
            if (invalid.Contains(ch))
                continue;

            sb.Append(ch);
        }

        return sb.ToString();
    }
}
