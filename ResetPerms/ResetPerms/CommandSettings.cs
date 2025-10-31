using Spectre.Console.Cli;
using System.ComponentModel;
using Spectre.Console;

namespace ResetPerms;

public sealed class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-s|--scriptsDirectory")]
    [Description("Path to the scripts directory")]
    public string? ScriptsDirectory { get; init; }
    
    public override ValidationResult Validate()
    {
        if (string.IsNullOrWhiteSpace(ScriptsDirectory))
        {
            return ValidationResult.Error("The --scriptsDirectory path must be provided.");
        }
        return base.Validate();
    }
}
