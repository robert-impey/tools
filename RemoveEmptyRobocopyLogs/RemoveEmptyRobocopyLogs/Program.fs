open System.IO
open FolderManager
open Microsoft.Extensions.FileSystemGlobbing

printfn "Looking for Robocopy Log Files"

let folderManager = FolderManager.GetFolderManager()

let synchLogsDirectory =
    Path.Combine(folderManager.GetLogsFolder(), "synch")

printfn "Synch logs directory: %s" synchLogsDirectory

let synchLogsDirMessage =
    if Directory.Exists(synchLogsDirectory) then
        $"synch logs directory %s{synchLogsDirectory} exists"
    else
        $"Synch logs directory %s{synchLogsDirectory} does not exist"

printfn $"%s{synchLogsDirMessage}"

let matcher = Matcher()
matcher.AddIncludePatterns(seq { "*.log"});
let matchingFiles = matcher.GetResultsInFullPath(synchLogsDirectory);

printfn "There are %d log files" (Seq.length matchingFiles)
