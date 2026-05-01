using System.ComponentModel;
using Spectre.Console;
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

    public override ValidationResult Validate()
    {
        if (string.IsNullOrWhiteSpace(Directory))
            return ValidationResult.Error("Directory is required.");

        if (!System.IO.Directory.Exists(Directory))
            return ValidationResult.Error($"Directory '{Directory}' does not exist.");

        // Optional: if logs-dir is provided, ensure it exists (or create it)
        if (!string.IsNullOrEmpty(LogsDirectory) && !System.IO.Directory.Exists(LogsDirectory))
            return ValidationResult.Error($"Logs directory '{LogsDirectory}' does not exist.");

        return ValidationResult.Success();
    }
}
