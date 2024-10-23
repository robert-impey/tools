open System.CommandLine
open FolderManager
open Microsoft.Extensions.Logging
open System.IO
open System

let async generateSynchWindowsConfigScript (logger : ILogger<FolderManager>) (filesFile: string) (synchScript: string) =
    logger.LogInformation $"Files file - {filesFile}; Synch script {synchScript}"  
    let outputFile = new StreamWriter(synchScript, false);

    outputFile.WriteLineAsync("# AUTOGEN'D - DO NOT EDIT!\n"
                            + $"# Written {DateTime.UtcNow:u}\n");

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
    
        if dryRun then
            logger.LogInformation "DRY RUN!"

        let folderManager = FolderManager.GetFolderManager(logger)

        let commonFilesFile = Path.Join(folderManager.GetCommonLocalScriptsFolder(), "synch", "config-Windows", "files.txt")
        logger.LogInformation commonFilesFile

        if File.Exists commonFilesFile then
            let scriptPath = Path.Join(folderManager.GetAutogenFolder(), "synch", "config-Windows.ps1")
            logger.LogInformation scriptPath

            if File.Exists scriptPath then
                logger.LogInformation "Deleting existing autogen'd script"
                File.Delete scriptPath

            generateSynchWindowsConfigScript logger commonFilesFile scriptPath
        else
            logger.LogInformation "Deleting existing autogen'd script"

    rootCommand.SetHandler(handler, dryRunOption)

    rootCommand.Invoke(args)
