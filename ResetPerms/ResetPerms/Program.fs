open System.IO
open System
open System.Diagnostics

type LocationFolder = {
    Location : string
    Folder : string
}

let findLatestErrFileFromTheLastDay logsDirectory =
    let logsDirectoryInfo = System.IO.DirectoryInfo(logsDirectory)

    let errFiles = logsDirectoryInfo.GetFiles("*.err")

    let lastDayFilter = fun (fi: FileInfo) ->
        fi.LastAccessTimeUtc > DateTime.UtcNow.AddDays(-1)

    let errFilesInLastDay = errFiles |> Array.where(lastDayFilter)

    let hasContentFilter = fun (fi: FileInfo) ->
        fi.Length > 0

    let errFilesWithContent = errFilesInLastDay |> Array.where (hasContentFilter)

    match errFilesWithContent.Length with
    | 0 -> None
    | 1 -> Some errFilesInLastDay.[0]
    | _ ->
        let sortedByDate = 
            errFilesInLastDay |> Array.sortBy (fun (fi: FileInfo) -> fi.LastAccessTimeUtc)
           
        Some sortedByDate.[0]

let extractFailedFile (file: string) =
    let prefix = "rsync: [sender] send_files failed to open \""
    let postfix = "\": Permission denied (13)"
    if file.StartsWith(prefix) && file.EndsWith(postfix) then
        let startIndex = prefix.Length
        let endIndex = file.Length - (prefix.Length + postfix.Length)
        Some (file.Substring(startIndex, endIndex))
    else
        None

let convertCygwinToWindows (file : string) =
    let prefix = "/cygdrive/"
    if file.StartsWith(prefix) then
        let startIndex = prefix.Length
        let shortened = file.Substring(startIndex)

        if shortened.[1] = '/' then
            Some (shortened.Substring(0, 1) + ":" + shortened.Substring(1))
        else
            None
    else
        None

let extractFailedFiles (errFile: FileInfo) =
    File.ReadAllLines(errFile.FullName)
        |> Array.choose (extractFailedFile) 
        |> Array.choose (convertCygwinToWindows)

let failedFileToLocationFolder (locations : string list) (failedFile: string) =
    let rec findLocation (locations' : string list) =
        match locations' with
        | [] -> None
        | x :: xs ->
            if failedFile.StartsWith(x) then
                Some x
            else
                findLocation xs
    
    let location = findLocation locations
    match location with
    | Some location' -> 
        let startIndex = location'.Length + 1
        let pathInLocation = failedFile.Substring(startIndex)

        let endIndex = pathInLocation.IndexOf('/')

        if endIndex < 0 then 
            None
        else
            let folder = pathInLocation.Substring(0, endIndex)
            Some ( {
                Location = location'
                Folder = folder
            })
    | None -> None

let locationFolderToBatch (resetScriptDirectory : string) (locationFolder: LocationFolder) =

    let location = locationFolder.Location.Replace('/', '_')

    let indexOfColon = location.IndexOf(':')

    let locationWithoutColon =
        if indexOfColon < 0 then location
        else
            location.Substring(0, indexOfColon) + location.Substring(indexOfColon + 1)
    let script = locationWithoutColon + "_" + locationFolder.Folder + ".bat"
    Path.Combine(resetScriptDirectory, script)

let searchLogsDirectory 
    (locations : string list) (resetScriptDirectory : string) (logsDirectory : string) =
    let latestErrFile = findLatestErrFileFromTheLastDay logsDirectory

    match latestErrFile with
    | Some latestErrFile' -> 
        let failedLocationFolders = 
            extractFailedFiles latestErrFile' 
            |> Array.choose (failedFileToLocationFolder locations)
            |> Array.distinct

        for failed in failedLocationFolders do
            let batch = locationFolderToBatch resetScriptDirectory failed
            if (File.Exists(batch)) then
                printfn "Executing %s" batch
                let batchProcess = Process.Start(batch)
                batchProcess.WaitForExit()
                batchProcess.Close()
            else
                printfn "Skipping %s" batch
    | None -> printfn "No error files found"

[<EntryPoint>]
let main(args) = 
    match args.Length with
    | 0 -> printfn "Tell me what to do!"
    | 1 -> 
        let logsDirectory = args.[0]

        let locations = [
            "c:/Users/rober/Dropbox"
            "c:/Users/rober/OneDrive"
            "x:"
            "z:"
        ]

        let resetScriptDirectory = 
            Path.Combine(
                System.Environment.GetEnvironmentVariable("USERPROFILE"),
                "autogen",
                "reset-perms")

        searchLogsDirectory locations resetScriptDirectory logsDirectory
    | _ -> printfn "Unrecognised input!"
    0