using System.Diagnostics;
using System.Runtime.InteropServices;

// See https://codepedia.info/dotnet-core-to-detect-operating-system-os-platform/
var isLinux = RuntimeInformation.IsOSPlatform(OSPlatform.Linux);
var isWindows = RuntimeInformation.IsOSPlatform(OSPlatform.Windows);

var homeDir = isWindows ? Environment.GetEnvironmentVariable("USERPROFILE") : Environment.GetEnvironmentVariable("HOME");

var localScriptsDir = isWindows ? Environment.GetEnvironmentVariable("LOCAL_SCRIPTS") : Path.Join(homeDir, "local-scripts");

var machineLocalScriptsDir =
    Path.Join(localScriptsDir, Environment.MachineName);

var userLocalScriptsDir =
    Path.Join(machineLocalScriptsDir, Environment.UserName);

const string stayDeletedFolder = "staydeleted";
const string nightly = "nightly.txt";
string? foundNightly = null;
var userNightly = Path.Join(userLocalScriptsDir, stayDeletedFolder, nightly);

if (File.Exists(userNightly))
{
    foundNightly = userNightly;
}
else
{
    var machineNightly = Path.Join(machineLocalScriptsDir, stayDeletedFolder, nightly);
    if (File.Exists(machineNightly))
    {
        foundNightly = machineNightly;
    }
}

if (foundNightly is null)
{
    Console.WriteLine("No nightly text file found!");
}
else
{
    Console.WriteLine($"Found {foundNightly}");

    // Where is the staydeleted executable?
    string? executablesDir = null;

    if (isWindows)
    {
        executablesDir = Environment.GetEnvironmentVariable("EXECUTABLES");
    }

    if (isLinux)
    {
        executablesDir = Path.Join(Environment.GetEnvironmentVariable("HOME"), "executables", "Linux", "prod", "x64");
    }

    string? stayDeletedExe = null;

    var exeSearch = Path.Join(executablesDir, "staydeleted");

    if (isWindows)
    {
        exeSearch = $"{exeSearch}.exe";
    }

    if (File.Exists(exeSearch))
    {
        stayDeletedExe = exeSearch;
    }

    // If present, invoke it with the found nightly.txt file.    
    if (stayDeletedExe is null)
    {
        Console.WriteLine("Stay Deleted executable not found!");
    }
    else
    {
        Console.WriteLine($"Running {stayDeletedExe} sweep {foundNightly}");

        var logsDir = Path.Join(homeDir, "logs");
        Process.Start(stayDeletedExe, new[] { "sweep", foundNightly, "--logs", logsDir, "--repeats", "1" });
    }
}
