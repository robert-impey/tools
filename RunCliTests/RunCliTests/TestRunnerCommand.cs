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
            if (cancellationToken.IsCancellationRequested)
            {
                AnsiConsole.MarkupLine("\n[yellow]Operation cancelled.[/]");
                return 1;
            }

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
                        if (await RunTestAsync(testFile, testDataDir, settings.Verbose, testType, programDir, cancellationToken))
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

    private async Task<bool> RunTestAsync(string testFile, string testDataDir, bool verbose, string testType, string programDir, CancellationToken cancellationToken)
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
        var (commandOutput, exitCode) = await ExecuteCommandAsync(command, programDir, testType, cancellationToken);

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

    private async Task<(string output, int exitCode)> ExecuteCommandAsync(string command, string workingDirectory, string outputType, CancellationToken cancellationToken)
    {
        var parts = command.Split(' ', 2, StringSplitOptions.RemoveEmptyEntries);
        var executable = parts[0];
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
            },
            EnableRaisingEvents = true // Allows events to fire
        };

        try
        {
            process.Start();

            // Asynchronously wait for the process to exit, passing the token
            await process.WaitForExitAsync(cancellationToken);

            // Note: If the token is cancelled, the process still runs but WaitForExitAsync throws.
            // We read the output regardless, as the process might have finished just before cancellation.

            // Read output asynchronously
            var outputTask = outputType == "out"
                ? process.StandardOutput.ReadToEndAsync()
                : process.StandardError.ReadToEndAsync();

            var output = await outputTask;

            return (output.TrimEnd('\r', '\n'), process.ExitCode);
        }
        catch (TaskCanceledException)
        {
            // If the user cancels the operation, kill the running process.
            if (!process.HasExited)
            {
                process.Kill();
            }
            throw; // Re-throw the exception to be caught by the calling function
        }
    }

    private void PrintSeparator(char character, int repetitions)
    {
        var line = new string(character, repetitions);
        AnsiConsole.WriteLine($"\n{line}\n");
    }
}
