open System
open System.ComponentModel
open System.IO
open System.Linq
open FolderManager
open Microsoft.Extensions.Logging
open Spectre.Console.Cli

let fileHasShebang (fileName: string) =
    use reader = new StreamReader(fileName)
    let firstLine = reader.ReadLine()

    firstLine.StartsWith("#!")

let findFilesWithShebang (scriptsDir: string) =
    let matchingFiles =
        Directory.EnumerateFiles(scriptsDir, "*", SearchOption.AllDirectories)

    matchingFiles |> Seq.filter fileHasShebang


type CliSettings() =
    inherit CommandSettings()

    [<CommandOption("-s|--scriptsDirectory")>]
    [<Description("Path to the scripts directory")>]
    member val ScriptsDirectory: string = "" with get, set
    
    [<CommandOption("--logged")>]
    [<Description("Logged or not")>]
    member val Logged: bool = false with get, set
    
    [<CommandOption("-l|--logsDirectory")>]
    [<Description("Path to the logs directory")>]
    member val LogsDirectory: string = "" with get, set

type DefaultCommand() =
    inherit Command<CliSettings>()

    override _.Execute (context: CommandContext, settings: CliSettings): int =
        let logger =
            if settings.Logged then
                if String.IsNullOrWhiteSpace(settings.LogsDirectory) then
                    raise (ArgumentNullException(settings.LogsDirectory))
                    
                LogsFileFinder.GetLogger<DefaultCommand>(settings.LogsDirectory, "ResetPerms")
            else
                use loggerFactory =
                    LoggerFactory.Create(fun builder ->
                        builder.AddConsole() |> ignore)
                loggerFactory.CreateLogger<DefaultCommand>()

        let filesWithShebang = findFilesWithShebang settings.ScriptsDirectory

        logger.LogInformation $"Found {filesWithShebang.Count()} files with shebangs"

        for file in filesWithShebang do
            logger.LogInformation $"File with shebang: {file}"

            File.SetUnixFileMode(
                file,
                UnixFileMode.UserRead
                ||| UnixFileMode.UserWrite
                ||| UnixFileMode.UserExecute
                ||| UnixFileMode.GroupRead
                ||| UnixFileMode.GroupExecute
                ||| UnixFileMode.OtherRead
                ||| UnixFileMode.OtherExecute
            )
        0

[<EntryPoint>]
let main args =
    let app = CommandApp<DefaultCommand>()
    app.Configure(fun config ->
        config.SetApplicationName("Reset Perms") |> ignore
    )
    app.Run(args)
