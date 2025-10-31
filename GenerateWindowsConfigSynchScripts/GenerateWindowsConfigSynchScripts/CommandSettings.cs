using Spectre.Console.Cli;
using System.ComponentModel;

namespace GenerateWindowsConfigSynchScripts;

public sealed class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-a|--autogen")]
    [Description("Autogen directory")]
    public string? Autogen { get; init; }

    [CommandOption("-c|--script")]
    [Description("Name of script to generate")]
    public string? Script { get; init; }

    [CommandOption("-f|--files <FILES>")]
    [Description("File with source, destination, and a list of files to be synch'd")]
    public string? Files { get; init; }
}
