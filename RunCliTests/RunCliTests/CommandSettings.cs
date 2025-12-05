using Spectre.Console.Cli;

namespace RunCliTests;

public sealed class CommandSettings : Spectre.Console.Cli.CommandSettings
{
    [CommandOption("-d|--directory <DIRECTORY>")]
    [System.ComponentModel.DefaultValue(".")]
    public string Directory { get; init; } = ".";

    [CommandOption("-v|--verbose")]
    [System.ComponentModel.DefaultValue(false)]
    public bool Verbose { get; init; }

    [CommandOption("-t|--test-data-directory <DIRECTORY>")]
    [System.ComponentModel.DefaultValue(".")]
    public string TestDataDirectory { get; init; } = ".";
}
