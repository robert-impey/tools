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
    let matchingFiles =
        Directory.EnumerateFiles(scriptsDir, "*", SearchOption.AllDirectories)

    matchingFiles |> Seq.filter fileHasShebang

[<EntryPoint>]
let main args =
    let rootCommand = RootCommand("Reset permissions for files that have a shebang")

    let handler () =
        let logger =
            LoggerFactory
                .Create(fun builder ->
                    builder.ClearProviders() |> ignore
                    builder.AddConsole() |> ignore)
                .CreateLogger<FolderManager>()

        let folderManager = FolderManager.GetFolderManager(logger)

        let filesWithShebang = findFilesWithShebang (folderManager.GetLocalScriptsFolder())

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

    rootCommand.SetHandler(handler)

    rootCommand.Invoke(args)
