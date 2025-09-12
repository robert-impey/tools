using Spectre.Console.Cli;
using System.ComponentModel;

namespace GenerateWindowsConfigSynchScripts;

public class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-a|--autogen")]
    [Description("Autogen directory")]
    public string? Autogen { get; set; }

    [CommandOption("-f|--files <FILES>")]
    [Description("File with a list of files to be synch'd")]
    public string? Files { get; set; }

    [CommandOption("-s|--source <SOURCE>")]
    [Description("Source directory")]
    public string? Source { get; set; }

    [CommandOption("-d|--destination <DEST>")]
    [Description("Destination directory")]
    public string? Destination { get; set; }

    [CommandOption("--logged")]
    [Description("Logged")]
    public bool Logged { get; set; }

    [CommandOption("-l|--logs <LOGS>")]
    [Description("Logs directory")]
    public string? LogsDirectory { get; set; }
}
