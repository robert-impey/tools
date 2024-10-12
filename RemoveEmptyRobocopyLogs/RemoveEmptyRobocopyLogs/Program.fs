open System.IO
open FolderManager
open Microsoft.Extensions.FileSystemGlobbing

open RobocopyLogs

printfn "Looking for Robocopy Log Files"

let synchLogsDirectory =
    Path.Combine(FolderManager.LogsFolder, "synch")

printfn "Synch logs directory: %s" synchLogsDirectory

let synchLogsDirMessage =
    if Directory.Exists(synchLogsDirectory) then
        $"synch logs directory %s{synchLogsDirectory} exists"
    else
        $"Synch logs directory %s{synchLogsDirectory} does not exist"

printfn $"%s{synchLogsDirMessage}"

let matcher = Matcher()
matcher.AddIncludePatterns(seq { "*.robocopy-synch.log"})
let matchingFiles = matcher.GetResultsInFullPath(synchLogsDirectory)

printfn "There are %d log files" (Seq.length matchingFiles)

for logFile in matchingFiles do
    if fileHasCopies logFile then
        printfn "%s has copies - keeping" logFile
    else
        printfn "%s has no copies - deleting" logFile
        File.Delete(logFile)
