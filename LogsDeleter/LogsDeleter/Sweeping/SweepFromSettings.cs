using System.ComponentModel;
using Spectre.Console.Cli;

namespace LogsDeleter.Sweeping;

public class SweepFromSettings : CommandSettings
{
    [CommandOption("-t|--tool <TOOL>")]
    [Description("Tool to sweep")]
    public string Tool { get; set; } = string.Empty;
}
