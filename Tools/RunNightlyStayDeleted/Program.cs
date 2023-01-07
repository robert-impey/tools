using System.Diagnostics;
using System.Runtime.InteropServices;

var machineLocalScriptsDir =
    Path.Join(Environment.GetEnvironmentVariable("HOME"), "local-scripts", Environment.MachineName);

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
    var executablesDir = Path.Join(Environment.GetEnvironmentVariable("HOME"), "executables");
    
    string? stayDeletedExe = null;
    
    // See https://codepedia.info/dotnet-core-to-detect-operating-system-os-platform/
    var isLinux = RuntimeInformation.IsOSPlatform(OSPlatform.Linux);
    if (isLinux)
    {
        var linuxExe = Path.Join(executablesDir, "Linux", "prod", "x64", "staydeleted");

        if (File.Exists(linuxExe))
        {
            stayDeletedExe = linuxExe;
        }
    }

    // If present, invoke it with the found nightly.txt file.    
    if (stayDeletedExe is null)
    {
        Console.WriteLine("Stay Deleted executable not found!");
    }
    else
    {
        Console.WriteLine($"Running {stayDeletedExe} sweep {foundNightly}");

        var logsDir = Path.Join(Environment.GetEnvironmentVariable("HOME"), "logs");
        Process.Start(stayDeletedExe, new[] { "sweep", foundNightly, "--logs", logsDir, "--repeats", "1" });
    }
}
