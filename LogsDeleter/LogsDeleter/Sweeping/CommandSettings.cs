using System.ComponentModel;
using Spectre.Console.Cli;

namespace LogsDeleter.Sweeping;

public class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-l|--logsDirectory <PATH>")]
    [Description("The logs directory")]
    public string LogsDirectory { get; set; } = string.Empty;

    [CommandOption("-d|--days <DAYS>")]
    [Description("Days ago for cut off")]
    [DefaultValue(30)]
    public int Days { get; set; }

    [CommandOption("-e|--deleteEmpty")]
    [Description("Delete empty log files")]
    public bool DeleteEmpty { get; set; }

    [CommandOption("-v|--verbose")]
    [Description("Verbose output")]
    public bool Verbose { get; set; }
}
