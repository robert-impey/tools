using System.Diagnostics;
using Tools;

var logsDeleter = ExecutablesHelper.FindExecutable("logs-deleter");

var logsDir = WellKnownFoldersHelper.GetLogsDir();
Console.WriteLine($"Logs Directory: {logsDir}");

var ldLogsDir = Path.Join(logsDir, "logs-deleter");
Console.WriteLine($"LD Logs Directory: {ldLogsDir}");

var logTime = DateTimeOffset.UtcNow.ToString("yyyy-MM-dd_HH_mm_ss");
Console.WriteLine($"Log Time: {logTime}");

var logFile = Path.Join(ldLogsDir, $"{logTime}.log");
Console.WriteLine($"Log File: {logFile}");

var errFile = Path.Join(ldLogsDir, $"{logTime}.err");
Console.WriteLine($"Err Time: {errFile}");

using var logWriter = File.CreateText(logFile);
using var errWriter = File.CreateText(errFile);

var processStartInfo = new ProcessStartInfo();
processStartInfo.CreateNoWindow = true;
processStartInfo.RedirectStandardOutput = true;
processStartInfo.RedirectStandardInput = true;
processStartInfo.UseShellExecute = false;
processStartInfo.Arguments = "sweepAll";
processStartInfo.FileName = logsDeleter;

var process = new Process();
process.StartInfo = processStartInfo;
// enable raising events because Process does not raise events by default
process.EnableRaisingEvents = true;
// attach the event handler for OutputDataReceived before starting the process
process.OutputDataReceived += delegate(object _, DataReceivedEventArgs e)
{
    // append the new data to the data already read-in
    logWriter.WriteLine(e.Data);
};

process.ErrorDataReceived += delegate(object _, DataReceivedEventArgs e)
{
    // append the new data to the data already read-in
    errWriter.WriteLine(e.Data);
};

// start the process
// then begin asynchronously reading the output
// then wait for the process to exit
// then cancel asynchronously reading the output
process.Start();
process.BeginOutputReadLine();
process.WaitForExit();
process.CancelOutputRead();
