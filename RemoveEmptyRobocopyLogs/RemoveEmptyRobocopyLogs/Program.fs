// For more information see https://aka.ms/fsharp-console-apps

open System.IO

printfn "Looking for Robocopy Log Files"

let folderManager = FolderManager.FolderManager.GetFolderManager()

let synchLogsDirectory =
    Path.Combine(folderManager.GetLogsFolder(), "synch")

printfn "Synch logs directory: %s" synchLogsDirectory

let synchLogsDirMessage =
    if Directory.Exists(synchLogsDirectory) then
        $"synch logs directory %s{synchLogsDirectory} exists"
    else
        $"Synch logs directory %s{synchLogsDirectory} does not exist"

printfn $"%s{synchLogsDirMessage}"
