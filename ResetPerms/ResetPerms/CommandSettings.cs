using Spectre.Console.Cli;
using System.ComponentModel;

namespace ResetPerms;

public class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-s|--scriptsDirectory")]
    [Description("Path to the scripts directory")]
    public string ScriptsDirectory { get; set; } = string.Empty;

    [CommandOption("--logged")]
    [Description("Logged or not")]
    public bool Logged { get; set; } = false;

    [CommandOption("-l|--logsDirectory")]
    [Description("Path to the logs directory")]
    public string? LogsDirectory { get; set; } = null;
}
