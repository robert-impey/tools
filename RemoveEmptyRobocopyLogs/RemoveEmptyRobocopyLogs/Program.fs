open System.IO
open FolderManager
open Microsoft.Extensions.FileSystemGlobbing
open System.Text.RegularExpressions

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

//let line = "    Files :      5117         0      5117         0         0         0"
//let line = "    Files :      5117         123      5117         0         0         0"

let fileHasCopies (fileName: string) =
    let mutable foundCopiesLine = false
    for line in File.ReadAllLines fileName do 
        let lineRegex = new Regex("\s*Files :\s+\d+\s+(\d+)\s+")

        let matches = lineRegex.Matches(line)

        if matches.Count > 0 then
            if "0" <> matches[0].Groups[1].Value then
                foundCopiesLine <- true

    foundCopiesLine

for logFile in matchingFiles do
    if fileHasCopies logFile then
        printfn "%s has copies - keeping" logFile
    else
        printfn "%s has no copies - deleting" logFile
        File.Delete(logFile)
