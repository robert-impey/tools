open System.CommandLine
open FolderManager
open Microsoft.Extensions.Logging
open System.IO
open System
open System.Text

let generateSynchWindowsConfigScript (logger : ILogger<FolderManager>) (filesFile: string) (synchScript: string) (homeFolder: string) (commonWindowsConfigFolder: string)=
    async {
        logger.LogInformation $"Writing to {synchScript}"
        use outputFile = File.OpenWrite(synchScript)

        do! outputFile.AsyncWrite("# AUTOGEN'D - DO NOT EDIT!\n" |> Encoding.ASCII.GetBytes)
        do! outputFile.AsyncWrite($"# Written {DateTime.UtcNow:u}\n\n" |> Encoding.ASCII.GetBytes)

        let writeScriptLine (file: string) =
            Async.RunSynchronously(outputFile.AsyncWrite($"ROBOCOPY {homeFolder} {commonWindowsConfigFolder} /xo {file}\n" |> Encoding.ASCII.GetBytes))
            Async.RunSynchronously(outputFile.AsyncWrite($"ROBOCOPY {commonWindowsConfigFolder} {homeFolder} /xo {file}\n\n" |> Encoding.ASCII.GetBytes))

        logger.LogInformation $"Reading {filesFile}"
        File.ReadAllLines filesFile 
            |> Seq.iter writeScriptLine
    }

[<EntryPoint>]
let main args =
    let loggedOption =
        Option<bool>("--logged", (fun () -> false))

    let rootCommand =
        RootCommand("Generate Windows Config Synch Scripts")

    rootCommand.AddOption(loggedOption)
    
    let handler (logged: bool) =
        let logger =
            if logged then
                LogsFileFinder.GetLogger<FolderManager>("synch", "GenerateWindowsConfigSynchScripts")
            else
                LoggerFactory.Create(fun builder ->
                    builder.ClearProviders() |> ignore
                    builder.AddConsole() |> ignore).CreateLogger<FolderManager>()

        let folderManager = FolderManager.GetFolderManager(logger)

        let commonFilesFile = Path.Join(folderManager.GetCommonLocalScriptsFolder(), "synch", "config-Windows", "files.txt")
        logger.LogInformation $"Files file - {commonFilesFile}"  

        if File.Exists commonFilesFile then
            let synchAutogen = Path.Join(folderManager.GetAutogenFolder(), "synch")

            if not (Directory.Exists(synchAutogen)) then
                Directory.CreateDirectory(synchAutogen) |> ignore

            let scriptPath = Path.Join(synchAutogen, "config-Windows.ps1")
            logger.LogInformation $"Script path - {scriptPath}"  

            if File.Exists scriptPath then
                logger.LogInformation "Deleting existing autogen'd script"
                File.Delete scriptPath

            let configFolder = Path.Join(FolderManager.ConfigFolder, "_Common", "Windows")

            if Directory.Exists(configFolder) then
                logger.LogInformation $"Common Windows config folder - {configFolder}"
                Async.RunSynchronously(generateSynchWindowsConfigScript logger commonFilesFile scriptPath FolderManager.HomeFolder configFolder)
            else
                logger.LogError $"Common Windows config folder does not exist - {configFolder}"
        else
            logger.LogInformation "Deleting existing autogen'd script"

    rootCommand.SetHandler(handler, loggedOption)

    rootCommand.Invoke(args)
