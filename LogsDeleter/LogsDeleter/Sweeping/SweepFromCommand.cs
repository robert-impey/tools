using Spectre.Console.Cli;

namespace LogsDeleter.Sweeping;

public class SweepFromCommand : Command<SweepFromSettings>
{
    private readonly LogsDeleter _logsDeleter;

    public SweepFromCommand(LogsDeleter logsDeleter)
    {
        _logsDeleter = logsDeleter;
    }

    protected override int Execute(CommandContext context, SweepFromSettings settings, CancellationToken cancellationToken)
    {
        if (string.IsNullOrEmpty(settings.Tool))
        {
            throw new Exception("Tool not set");
        }

        if (string.IsNullOrEmpty(settings.LogsDirectory))
        {
            throw new Exception("LogsDirectory not set");
        }

        string toolPath = Path.Combine(settings.LogsDirectory, settings.Tool);

        if (!Directory.Exists(toolPath))
        {
            throw new DirectoryNotFoundException($"Tool path not found: {toolPath}");
        }

        _logsDeleter.DeleteFrom(toolPath, settings.Days, settings.DeleteEmpty, settings.Verbose);

        if (settings.Verbose) Console.WriteLine("Success");
        return 0;
    }
}
