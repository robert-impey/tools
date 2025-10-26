using Spectre.Console;
using Spectre.Console.Cli;
using System.ComponentModel;

namespace RemoveEmptyRobocopyLogs;

public sealed class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-l|--logsDirectory <PATH>")]
    [Description("Path to the logs directory")]
    public string LogsDirectory { get; set; } = string.Empty;

    // Optional: Add validation logic
    public override ValidationResult Validate()
    {
        if (string.IsNullOrWhiteSpace(LogsDirectory))
        {
            return ValidationResult.Error("The --logsDirectory path must be provided.");
        }
        return base.Validate();
    }
}
