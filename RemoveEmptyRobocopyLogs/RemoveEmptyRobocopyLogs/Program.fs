open System.IO
open Microsoft.Extensions.FileSystemGlobbing
open Spectre.Console.Cli
open System.ComponentModel
open Microsoft.Extensions.Logging

open RobocopyLogs

type CliSettings() =
    inherit CommandSettings()

    [<CommandOption("-l|--logsDirectory")>]
    [<Description("Path to the logs directory")>]
    member val LogsDirectory: string = "" with get, set

type DefaultCommand() =
    inherit Command<CliSettings>()

    override _.Execute(_: CommandContext, settings: CliSettings) =
        let logger =
            use loggerFactory =
                LoggerFactory.Create(fun builder ->
                    builder.AddConsole() |> ignore)
            loggerFactory.CreateLogger<DefaultCommand>()
        
        logger.LogInformation "Looking for Robocopy Log Files"

        logger.LogInformation $"Logs directory: {settings.LogsDirectory}"

        if Directory.Exists(settings.LogsDirectory) then
            logger.LogInformation $"synch logs directory %s{settings.LogsDirectory} exists"
        else
            logger.LogInformation $"Synch logs directory %s{settings.LogsDirectory} does not exist"

        let matcher = Matcher()
        matcher.AddIncludePatterns(seq { "*.robocopy-synch.log"})
        let matchingFiles = matcher.GetResultsInFullPath(settings.LogsDirectory)

        logger.LogInformation $"There are {Seq.length matchingFiles} log files"

        for logFile in matchingFiles do
            if fileHasCopies logFile then
                logger.LogInformation $"{logFile} has copies - keeping"
            else
                logger.LogInformation $"{logFile} has no copies - deleting"
                File.Delete(logFile)
        0

[<EntryPoint>]
let main args =
    let app = CommandApp<DefaultCommand>()
    app.Configure(fun config ->
        config.SetApplicationName("RemoveEmptyRobocopyLogs") |> ignore
    )
    app.Run(args)
