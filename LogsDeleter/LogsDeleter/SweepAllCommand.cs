using Spectre.Console.Cli;

namespace LogsDeleter;

public class SweepAllCommand : Command<LogSettings>
{
    private readonly LogsDeleter _logsDeleter;

    public SweepAllCommand(LogsDeleter logsDeleter)
    {
        _logsDeleter = logsDeleter;
    }

    public override int Execute(CommandContext context, LogSettings settings, CancellationToken cancellationToken)
    {
        if (string.IsNullOrEmpty(settings.LogsDirectory))
        {
            throw new Exception("LogsDirectory not set");
        }

        const string toolName = "logs-deleter";
        string toolLogDir = Path.Combine(settings.LogsDirectory, toolName);

        if (!Directory.Exists(toolLogDir))
        {
            Directory.CreateDirectory(toolLogDir);
        }

        var subDirs = Directory.GetDirectories(settings.LogsDirectory);

        foreach (var subDir in subDirs)
        {
            _logsDeleter.DeleteFrom(subDir, settings.Days, settings.DeleteEmpty, settings.Verbose);
        }

        if (settings.Verbose) Console.WriteLine("Success");
        return 0;
    }


}
