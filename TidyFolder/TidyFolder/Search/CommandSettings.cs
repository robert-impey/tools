using System.ComponentModel;
using Spectre.Console.Cli;

namespace TidyFolder.Search;

public class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-l|--logs-dir <PATH>")]
    [Description("Directory where logs should be written")]
    public string LogsDirectory { get; set; } = string.Empty;

    [CommandArgument(0, "<directory>")]
    [Description("Search a single directory")]
    public string Directory { get; init; } = string.Empty;
}
