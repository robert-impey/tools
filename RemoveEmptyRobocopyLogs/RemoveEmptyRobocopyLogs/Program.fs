// For more information see https://aka.ms/fsharp-console-apps

open System.IO

printfn "Looking for Robocopy Log Files"

let logsDirectory =
    Path.Combine(System.Environment.GetEnvironmentVariable("HOME"), "logs", "synch")

printfn "Synch logs directory: %s" logsDirectory

let synchLogsDirMessage =
    if Directory.Exists(logsDirectory) then
        $"synch logs directory %s{logsDirectory} exists"
    else
        $"Synch logs directory %s{logsDirectory} does not exist"

printfn $"%s{synchLogsDirMessage}"
