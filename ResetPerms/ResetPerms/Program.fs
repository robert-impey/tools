open System.IO
open System
open System.CommandLine

let findDefaultLocationsFile () =
    let name =
        Path.Combine(
            Environment.GetEnvironmentVariable("USERPROFILE"),
            "Dropbox",
            "local-scripts",
            Environment.GetEnvironmentVariable("ComputerName"),
            "reset-perms",
            "locations.txt"
        )

    if (File.Exists(name)) then
        Some name
    else
        None

[<EntryPoint>]
let main (args) =
    let dryRunOption =
        Option<bool>("--dry-run", (fun () -> false))

    let locationsFileOption =
        Option<string option>("--locations-file", findDefaultLocationsFile)
    let rootCommand =
        RootCommand("Reset permissions for files that failed synch")

    rootCommand.AddOption(dryRunOption)
    rootCommand.AddOption(locationsFileOption)

    let handler (dryRun: bool) (locationsFile : string option) (errFile : string option) =
        match locationsFile with
        | None -> printfn "No locations file found in the default location - quitting"
        | Some locationsFile' ->
                printfn $"%b{dryRun} %s{locationsFile'}"

    rootCommand.SetHandler(handler, dryRunOption, locationsFileOption)

    rootCommand.Invoke(args)
