open System.CommandLine
open FolderManager
open Microsoft.Extensions.Logging

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
                LogsFileFinder.GetLogger<FolderManager>("synch", "GenerateWindowsConfigSynchScripts")
    
        let folderManager = FolderManager.GetFolderManager(logger)

        if dryRun then
            logger.LogInformation "DRY RUN!"
        else
            logger.LogInformation "NOT DRY RUN!"

    rootCommand.SetHandler(handler, dryRunOption)

    rootCommand.Invoke(args)
