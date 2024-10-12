open System.CommandLine
open System.Linq
open FolderManager
open Microsoft.Extensions.FileSystemGlobbing

let findFilesWithShebang (scriptsDir: string) =
    let matcher = Matcher()
    matcher.AddIncludePatterns(seq { "*"})
    let matchingFiles = matcher.GetResultsInFullPath(scriptsDir)
    matchingFiles

[<EntryPoint>]
let main (args) =
    let dryRunOption =
        Option<bool>("--dry-run", (fun () -> false))

    let rootCommand =
        RootCommand("Reset permissions for files that have a shebang")

    rootCommand.AddOption(dryRunOption)
    
    let handler (dryRun: bool) =
        let logger =
            if dryRun then
                let config = NLog.Config.LoggingConfiguration ()

                let logconsole = new NLog.Targets.ConsoleTarget("logconsole")

                config.AddRule(NLog.LogLevel.Info, NLog.LogLevel.Fatal, logconsole)

                NLog.LogManager.Configuration = config |> ignore
                
                NLog.LogManager.GetCurrentClassLogger()
            else
                LogsFileFinder.GetLogger("reset-perms", "ResetPerms")
    
        let folderManager = FolderManager.GetFolderManager(logger)
    
        let filesWithShebang = findFilesWithShebang (folderManager.GetLocalScriptsFolder ())
        
        if dryRun then
            logger.Info "DRY RUN!"
        
        logger.Info $"Found {filesWithShebang.Count()} files with shebangs"
        
        for file in filesWithShebang do
            logger.Info file

    rootCommand.SetHandler(handler, dryRunOption)

    rootCommand.Invoke(args)
