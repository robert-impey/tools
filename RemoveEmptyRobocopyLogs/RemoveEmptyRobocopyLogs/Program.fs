open System.IO
open Microsoft.Extensions.FileSystemGlobbing
open Spectre.Console
open Spectre.Console.Cli
open System.ComponentModel

open RobocopyLogs

type CliSettings() =
    inherit CommandSettings()

    [<CommandOption("-l|--logsDirectory")>]
    [<Description("Path to the logs directory")>]
    member val LogsDirectory: string = "" with get, set

type DefaultCommand() =
    inherit Command<CliSettings>()

    override _.Execute(context, settings) =
        printfn "Looking for Robocopy Log Files"

        AnsiConsole.MarkupLine($"[green]Logs directory:[/] {settings.LogsDirectory}")

        let synchLogsDirMessage =
            if Directory.Exists(settings.LogsDirectory) then
                $"synch logs directory %s{settings.LogsDirectory} exists"
            else
                $"Synch logs directory %s{settings.LogsDirectory} does not exist"

        printfn $"%s{synchLogsDirMessage}"

        let matcher = Matcher()
        matcher.AddIncludePatterns(seq { "*.robocopy-synch.log"})
        let matchingFiles = matcher.GetResultsInFullPath(settings.LogsDirectory)

        printfn "There are %d log files" (Seq.length matchingFiles)

        for logFile in matchingFiles do
            if fileHasCopies logFile then
                printfn "%s has copies - keeping" logFile
            else
                printfn "%s has no copies - deleting" logFile
                File.Delete(logFile)
        0

[<EntryPoint>]
let main args =
    let app = CommandApp<DefaultCommand>()
    app.Configure(fun config ->
        config.SetApplicationName("LogViewer") |> ignore
    )
    app.Run(args)
