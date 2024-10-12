open System.CommandLine
open System.IO
open System.Linq
open FolderManager
open Microsoft.Extensions.Logging

let fileHasShebang (fileName: string) =
    let firstLine = File.ReadAllLines fileName |> Seq.take 1 |> Seq.toArray
    
    match firstLine.Length with
    | 1 -> firstLine.[0].StartsWith("#!")
    | _ -> false

let findFilesWithShebang (scriptsDir: string) =
    let matchingFiles = Directory.EnumerateFiles(scriptsDir, "*", SearchOption.AllDirectories)
    matchingFiles |> Seq.filter fileHasShebang

[<EntryPoint>]
let main args =
    let dryRunOption =
        Option<bool>("--dry-run", (fun () -> false))

    let rootCommand =
        RootCommand("Reset permissions for files that have a shebang")

    rootCommand.AddOption(dryRunOption)
    
    let handler (dryRun: bool) =
        let logger =
            if dryRun then
                LoggerFactory.Create(fun builder ->
                    builder.ClearProviders() |> ignore
                    builder.AddConsole() |> ignore).CreateLogger<FolderManager>()
            else
                LogsFileFinder.GetLogger<FolderManager>("reset-perms", "ResetPerms")
    
        let folderManager = FolderManager.GetFolderManager(logger)
    
        let filesWithShebang = findFilesWithShebang (folderManager.GetLocalScriptsFolder ())
        
        if dryRun then
            logger.LogInformation "DRY RUN!"
        
        logger.LogInformation $"Found {filesWithShebang.Count()} files with shebangs"
        
        for file in filesWithShebang do
            logger.LogInformation file

    rootCommand.SetHandler(handler, dryRunOption)

    rootCommand.Invoke(args)
