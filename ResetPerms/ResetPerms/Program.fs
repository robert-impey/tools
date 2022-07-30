open System.IO
open System
open System.Diagnostics
open System.CommandLine

open ResetPerms.Cygwin

type LocationFolder = { Location: string; Folder: string }

let findLatestErrFile logsDirectory =
    let logsDirectoryInfo =
        System.IO.DirectoryInfo(logsDirectory)

    let errFiles =
        logsDirectoryInfo.GetFiles("*.err")

    let hasContentFilter =
        fun (fi: FileInfo) -> fi.Length > 0

    let errFilesWithContent =
        errFiles |> Array.where (hasContentFilter)

    match errFilesWithContent.Length with
    | 0 -> None
    | 1 -> Some errFiles.[0]
    | _ ->
        let sortedByDate =
            errFiles
            |> Array.sortByDescending (fun (fi: FileInfo) -> fi.LastAccessTimeUtc)

        Some sortedByDate.[0]

let extractFailedFiles (errFile: FileInfo) =
    File.ReadAllLines(errFile.FullName)
    |> Array.choose (extractFailedFile)
    |> Array.choose (convertCygwinToWindows)

let failedFileToLocationFolder (locations: string list) (failedFile: string) =
    let rec findLocation (locations': string list) =
        match locations' with
        | [] -> None
        | x :: xs ->
            if failedFile.StartsWith(x, StringComparison.InvariantCultureIgnoreCase) then
                Some x
            else
                findLocation xs

    let location = findLocation locations

    match location with
    | Some location' ->
        let startIndex = location'.Length + 1

        let pathInLocation =
            failedFile.Substring(startIndex)

        let endIndex = pathInLocation.IndexOf('\\')

        if endIndex < 0 then
            None
        else
            let folder =
                pathInLocation.Substring(0, endIndex)

            Some(
                { Location = location'
                  Folder = folder }
            )
    | None -> None

let locationFolderToBatch (resetScriptDirectory: string) (locationFolder: LocationFolder) =

    let location =
        locationFolder.Location.Replace('\\', '_')

    let indexOfColon = location.IndexOf(':')

    let locationWithoutColon =
        if indexOfColon < 0 then
            location
        else
            location.Substring(0, indexOfColon)
            + location.Substring(indexOfColon + 1)

    let script =
        locationWithoutColon
        + "_"
        + locationFolder.Folder
        + ".bat"

    Path.Combine(resetScriptDirectory, script)

let resetPermsFolder (dryRun: bool) (locationFolder: LocationFolder) (resetScriptDirectory: string) =
    let batch = locationFolderToBatch resetScriptDirectory locationFolder

    if (File.Exists(batch)) then
        printfn "Executing %s" batch

        if dryRun then
            printfn "DRY RUN (execution skipped)"
        else
            let batchProcess = Process.Start(batch)
            batchProcess.WaitForExit()
            batchProcess.Close()
    else
        printfn "Skipping %s" batch

let searchLogsDirectory (dryRun: bool) (locations: string list) (resetScriptDirectory: string) (logsDirectory: string) =
    let latestErrFile =
        findLatestErrFile logsDirectory

    match latestErrFile with
    | Some latestErrFile' ->
        let failedLocationFolders =
            extractFailedFiles latestErrFile'
            |> Array.choose (failedFileToLocationFolder locations)
            |> Array.distinct

        for failed in failedLocationFolders do
            resetPermsFolder dryRun failed resetScriptDirectory
    | None -> printfn "No error files found"

let findDefaultLocationsFile () =
    let name =
        Path.Combine(
            System.Environment.GetEnvironmentVariable("USERPROFILE"),
            "Dropbox",
            "local-scripts",
            System.Environment.GetEnvironmentVariable("ComputerName"),
            "reset-perms",
            "locations.txt"
        )

    if (File.Exists(name)) then
        Some name
    else
        None

let readLocationsFile (locationsFile: string) =
    let f =
        fun s ->
            if String.IsNullOrWhiteSpace(s) then
                None
            else
                Some s

    File.ReadAllLines(locationsFile)
    |> Array.choose (f)
    |> Array.toList

[<EntryPoint>]
let main (args) =
    let dryRunOption =
        Option<bool>("--dry-run", (fun () -> false))

    let locationsFile =
        findDefaultLocationsFile ()

    let resetScriptDirectory =
        Path.Combine(System.Environment.GetEnvironmentVariable("USERPROFILE"), "autogen", "reset-perms")

    let logsDirectory =
        Path.Combine(System.Environment.GetEnvironmentVariable("USERPROFILE"), "logs", "synch")

    let rootCommand =
        RootCommand("Reset permissions for files that failed synch")

    rootCommand.AddOption(dryRunOption)

    rootCommand.SetHandler(
        fun (dryRun: bool) ->
            match locationsFile with
            | Some locationsFile' ->
                let locations =
                    readLocationsFile locationsFile'

                searchLogsDirectory dryRun locations resetScriptDirectory logsDirectory
            | None -> printfn "No locations file found"
        ,
        dryRunOption)

    rootCommand.Invoke(args)
