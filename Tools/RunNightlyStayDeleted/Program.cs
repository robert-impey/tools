using System.Diagnostics;
using Tools;

// See https://codepedia.info/dotnet-core-to-detect-operating-system-os-platform/
var isWindows = OsHelper.IsWindows();

var homeDir = OsHelper.GetHomeDir();

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

    var stayDeletedExe = ExecutablesHelper.FindExecutable("staydeleted");

    Console.WriteLine($"Stay Deleted Exe: {stayDeletedExe}");
    Console.WriteLine($"Nightly File: {foundNightly}");

    var logsDir = WellKnownFoldersHelper.GetLogsDir();
    Console.WriteLine($"Logs Directory: {logsDir}");

    var repeats = isWindows ? 8 : 1;
    Console.WriteLine($"Repeats: {repeats}");

    Process.Start(stayDeletedExe, new[] { "sweep", foundNightly, "--logs", logsDir, "--repeats", repeats.ToString() });
}
