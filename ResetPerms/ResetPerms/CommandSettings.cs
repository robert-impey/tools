using Spectre.Console.Cli;
using System.ComponentModel;
using Spectre.Console;

namespace ResetPerms;

public sealed class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-s|--scriptsDirectory")]
    [Description("Path to the scripts directory")]
    public string? ScriptsDirectory { get; init; }

    [CommandOption("--logged")]
    [Description("Logged or not")]
    [DefaultValue(false)]
    public bool Logged { get; init; }

    [CommandOption("-l|--logsDirectory")]
    [Description("Path to the logs directory")]
    public string? LogsDirectory { get; init; }
    
    public override ValidationResult Validate()
    {
        if (string.IsNullOrWhiteSpace(ScriptsDirectory))
        {
            return ValidationResult.Error("The --scriptsDirectory path must be provided.");
        }
        return base.Validate();
    }
}
