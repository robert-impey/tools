using System.Diagnostics;
using Spectre.Console;
using Spectre.Console.Cli;

namespace RunCliTests;

public sealed class TestRunnerCommand : AsyncCommand<CommandSettings>
{
    private const string TestsDirName = "tests";

    // --- Entry Point ---
    public override async Task<int> ExecuteAsync(CommandContext context, CommandSettings settings, CancellationToken cancellationToken)
    {
        // 1. Normalize and Absolutize Paths
        var dir = Path.GetFullPath(settings.Directory);
        var testDataDir = Path.GetFullPath(settings.TestDataDirectory);

        var successes = 0;
        var tests = 0;

        if (settings.Verbose)
        {
            AnsiConsole.WriteLine($"Searching {dir}");
        }

        // 2. Directory Traversal (Recursive search)
        // Includes the root directory and all non-hidden subdirectories
        var programDirs = Directory.EnumerateDirectories(dir, "*", SearchOption.AllDirectories)
            .Prepend(dir) // Include the starting directory itself
            .Where(d => !Path.GetFileName(d).StartsWith(".")) // Skip hidden directories
            .ToList();

        foreach (var programDir in programDirs)
        {
            var testDirPath = Path.Combine(programDir, TestsDirName);

            if (settings.Verbose)
            {
                PrintSeparator('+', 40);
                AnsiConsole.MarkupLine($"[yellow]Program directory:[/] [bold]{programDir}[/]");
            }

            if (Directory.Exists(testDirPath))
            {
                var testFiles = Directory.EnumerateFiles(testDirPath)
                    .OrderBy(f => Path.GetFileName(f)); // Sort for consistent order

                // Group and run .txt and .err tests
                foreach (var testFile in testFiles)
                {
                    var extension = Path.GetExtension(testFile).ToLowerInvariant();
                    string testType = string.Empty;

                    if (extension == ".txt")
                    {
                        testType = "out";
                    }
                    else if (extension == ".err")
                    {
                        testType = "err";
                    }

                    if (!string.IsNullOrEmpty(testType))
                    {
                        tests++;
                        if (RunTest(testFile, testDataDir, settings.Verbose, testType, programDir))
                        {
                            successes++;
                        }
                    }
                }
            }

            if (settings.Verbose)
            {
                PrintSeparator('+', 40);
            }
        }

        // 3. Final Summary
        if (tests == 0)
        {
            AnsiConsole.MarkupLine("[red]No tests found![/]");
            return 1; // Return non-zero for error/warning
        }
        else
        {
            var successRate = (double)successes / tests * 100.0;
            AnsiConsole.MarkupLine($"Success rate: [green]{successes}[/]/[bold]{tests}[/] ([green]{successRate:F1}%[/])");
        }

        return 0; // Success
    }

    // --- Helper Methods ---

    private bool RunTest(string testFile, string testDataDir, bool verbose, string testType, string programDir)
    {
        if (verbose) PrintSeparator('-', 40);

        var fileName = Path.GetFileName(testFile);
        AnsiConsole.Write($"Test file: {fileName}{(verbose ? "\n" : " ")}");

        // Read command and expected output
        var (command, expectedOutput) = ReadTestFile(testFile, testDataDir, verbose);

        if (verbose)
        {
            AnsiConsole.WriteLine("Test output (Expected):");
            PrintSeparator('.', 40);
            AnsiConsole.WriteLine(expectedOutput);
            PrintSeparator('.', 40);
        }

        // Execute Command
        var (commandOutput, exitCode) = ExecuteCommand(command, programDir, testType);

        if (verbose)
        {
            AnsiConsole.WriteLine("Command output (Actual):");
            PrintSeparator('.', 40);
            AnsiConsole.WriteLine(commandOutput);
            PrintSeparator('.', 40);
        }

        // Compare output
        var success = commandOutput.TrimEnd('\r', '\n') == expectedOutput.TrimEnd('\r', '\n');

        if (success)
        {
            AnsiConsole.MarkupLine("[green]OK[/]");
        }
        else
        {
            AnsiConsole.MarkupLine("[red]FAIL[/]");
        }

        if (verbose) PrintSeparator('-', 40);

        return success;
    }

    private (string command, string testOutput) ReadTestFile(string testFile, string testDataDir, bool verbose)
    {
        var lines = File.ReadAllLines(testFile);

        if (lines.Length < 3)
        {
            throw new InvalidOperationException($"Malformed test file: {testFile}. Must contain at least 3 lines.");
        }

        // Line 1: Command
        var command = lines[0].Trim();

        if (string.IsNullOrEmpty(command))
        {
            throw new InvalidOperationException($"Empty command in {testFile}!");
        }

        // Handle path substitution
        command = command.Replace("TEST_DATA_DIR", testDataDir);

        // Prepend "./" for relative paths if necessary (Approximation of Perl's behavior)
        if (!Path.IsPathRooted(command) && !command.StartsWith("./"))
        {
            command = $"./{command}";
        }

        if (verbose)
        {
            AnsiConsole.WriteLine($"Command: {command}");
        }

        // Line 2 is ignored

        // Remaining lines: Expected Test Output
        // Skip the first two lines and rejoin the rest
        var testOutput = string.Join(Environment.NewLine, lines.Skip(2));

        return (command, testOutput);
    }

    private (string output, int exitCode) ExecuteCommand(string command, string workingDirectory, string outputType)
    {
        var parts = command.Split(' ', 2, StringSplitOptions.RemoveEmptyEntries);
        var executable = parts.Length > 0 ? parts[0] : string.Empty;
        var arguments = parts.Length > 1 ? parts[1] : string.Empty;

        using var process = new Process
        {
            StartInfo = new ProcessStartInfo
            {
                FileName = executable,
                Arguments = arguments,
                WorkingDirectory = workingDirectory,
                RedirectStandardOutput = true,
                RedirectStandardError = true,
                UseShellExecute = false,
                CreateNoWindow = true
            }
        };

        process.Start();
        process.WaitForExit();

        var output = outputType == "out"
            ? process.StandardOutput.ReadToEnd()
            : process.StandardError.ReadToEnd();

        return (output.TrimEnd('\r', '\n'), process.ExitCode);
    }

    private void PrintSeparator(char character, int repetitions)
    {
        var line = new string(character, repetitions);
        AnsiConsole.WriteLine($"\n{line}\n");
    }
}
