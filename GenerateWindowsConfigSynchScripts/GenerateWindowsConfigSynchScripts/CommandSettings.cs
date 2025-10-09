using Spectre.Console.Cli;
using System.ComponentModel;

namespace GenerateWindowsConfigSynchScripts;

public class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-a|--autogen")]
    [Description("Autogen directory")]
    public string? Autogen { get; set; }

    [CommandOption("-c|--script")]
    [Description("Name of script to generate")]
    public string? Script { get; set; }

    [CommandOption("-f|--files <FILES>")]
    [Description("File with source, destination, and a list of files to be synch'd")]
    public string? Files { get; set; }

    [CommandOption("--logged")]
    [Description("Logged or not")]
    public bool Logged { get; set; }

    [CommandOption("-l|--logs <LOGS>")]
    [Description("Logs directory")]
    public string? LogsDirectory { get; set; }
}
