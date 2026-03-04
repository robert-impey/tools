using System.Diagnostics;
using Spectre.Console;
using Spectre.Console.Cli;

namespace RunCliTests;

public sealed class DefaultCommand : AsyncCommand<CommandSettings>
{
    private const string TestsDirName = "tests";

    public override async Task<int> ExecuteAsync(
        CommandContext context,
        CommandSettings settings,
        CancellationToken cancellationToken
    )
    {
        var dir = Path.GetFullPath(settings.Directory);
        var testDataDir = Path.GetFullPath(settings.TestDataDirectory);

        var successes = 0;
        var tests = 0;

        if (settings.Verbose)
        {
            AnsiConsole.WriteLine($"Searching {dir}");
        }

        var programDirs = EnumerateFilteredDirectories(dir);

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
                    .OrderBy(Path.GetFileName); // Sort for consistent order

                // Group and run .txt and .err tests
                foreach (var testFile in testFiles)
                {
                    var extension = Path.GetExtension(testFile).ToLowerInvariant();

                    var testType = extension switch
                    {
                        ".txt" => "out",
                        ".err" => "err",
                        _ => string.Empty
                    };

                    if (string.IsNullOrEmpty(testType))
                    {
                        continue;
                    }

                    tests++;
                    if (await RunTestAsync(testFile, testDataDir, settings.Verbose, testType, programDir,
                            cancellationToken))
                    {
                        successes++;
                    }
                }
            }

            if (settings.Verbose)
            {
                PrintSeparator('+', 40);
            }
        }

        if (tests == 0)
        {
            AnsiConsole.MarkupLine("[red]No tests found![/]");
            return 1; // Return non-zero for error/warning
        }

        var successRate = (double)successes / tests * 100.0;
        AnsiConsole.MarkupLine($"Success rate: [green]{successes}[/]/[bold]{tests}[/] ([green]{successRate:F1}%[/])");

        return 0; // Success
    }

    private static List<string> EnumerateFilteredDirectories(string rootDir)
    {
        var programDirs = new List<string> { rootDir }; // Start with the root directory

        // Use a queue for a breadth-first search (more memory efficient than pure recursion)
        var queue = new Queue<string>();
        queue.Enqueue(rootDir);

        while (queue.Count > 0)
        {
            string currentDir = queue.Dequeue();

            foreach (string subDir in Directory.EnumerateDirectories(currentDir))
            {
                string dirName = Path.GetFileName(subDir);

                // Apply the filter BEFORE adding to the list and before queuing for further search
                if (!dirName.StartsWith('.') && !subDir.EndsWith(".dSYM"))
                {
                    programDirs.Add(subDir);
                    queue.Enqueue(subDir); // Queue the valid directory for further searching
                }
                // If the filter fails, we simply skip it and do NOT queue it,
                // preventing the search from descending into it.
            }
        }

        return programDirs;
    }


    private async static Task<bool> RunTestAsync(
        string testFile,
        string testDataDir,
        bool verbose,
        string testType,
        string programDir,
        CancellationToken cancellationToken
    )
    {
        if (verbose)
        {
            PrintSeparator('-', 40);
        }

        var fileName = Path.GetFileName(testFile);
        AnsiConsole.Write($"Test file: {fileName}{(verbose ? "\n" : " ")}");

        var (command, expectedOutput) = ReadTestFile(testFile, testDataDir, verbose);

        if (verbose)
        {
            AnsiConsole.WriteLine("Test output (Expected):");
            PrintSeparator('.', 40);
            AnsiConsole.WriteLine(expectedOutput);
            PrintSeparator('.', 40);
        }

        var (commandOutput, exitCode) = await ExecuteCommandAsync(command, programDir, testType, cancellationToken);

        if (verbose)
        {
            AnsiConsole.WriteLine("Command output (Actual):");
            PrintSeparator('.', 40);
            AnsiConsole.WriteLine(commandOutput);
            PrintSeparator('.', 40);
            AnsiConsole.WriteLine($"Exit code: {exitCode}");
            PrintSeparator('.', 40);
        }

        var success = commandOutput.TrimEnd('\r', '\n') == expectedOutput.TrimEnd('\r', '\n');

        AnsiConsole.MarkupLine(success ? "[green]OK[/]" : "[red]FAIL[/]");

        if (verbose)
        {
            PrintSeparator('-', 40);
        }

        return success;
    }

    private static (string command, string testOutput) ReadTestFile(string testFile, string testDataDir, bool verbose)
    {
        var lines = File.ReadAllLines(testFile);

        if (lines.Length < 3)
        {
            throw new InvalidOperationException($"Malformed test file: {testFile}. Must contain at least 3 lines.");
        }

        var command = lines[0].Trim();

        if (string.IsNullOrEmpty(command))
        {
            throw new InvalidOperationException($"Empty command in {testFile}!");
        }

        // Handle path substitution
        command = command.Replace("TEST_DATA_DIR", testDataDir);

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

    private async static Task<(string output, int exitCode)> ExecuteCommandAsync(
        string command,
        string workingDirectory,
        string outputType,
        CancellationToken cancellationToken
    )
    {
        var parts = command.Split(' ', 2, StringSplitOptions.RemoveEmptyEntries);
        var executable = parts[0];
        var arguments = parts.Length > 1 ? parts[1] : string.Empty;

        var executablePath = Path.Join(workingDirectory, executable);

        using var process = new Process();
        process.StartInfo = new ProcessStartInfo
        {
            FileName = executablePath,
            Arguments = arguments,
            WorkingDirectory = workingDirectory,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            UseShellExecute = false,
            CreateNoWindow = true
        };
        process.EnableRaisingEvents = true; // Allows events to fire

        try
        {
            process.Start();

            // Asynchronously wait for the process to exit, passing the token
            await process.WaitForExitAsync(cancellationToken);

            // Note: If the token is cancelled, the process still runs but WaitForExitAsync throws.
            // We read the output regardless, as the process might have finished just before cancellation.

            // Read output asynchronously
            var outputTask = outputType == "out"
                ? process.StandardOutput.ReadToEndAsync(cancellationToken)
                : process.StandardError.ReadToEndAsync(cancellationToken);

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

            throw;
        }
    }

    private static void PrintSeparator(char character, int repetitions)
    {
        var line = new string(character, repetitions);
        AnsiConsole.WriteLine($"\n{line}\n");
    }
}
