open System.CommandLine

[<EntryPoint>]
let main (args) =
    let dryRunOption =
        Option<bool>("--dry-run", (fun () -> false))

    let rootCommand =
        RootCommand("Reset permissions for files that failed synch")

    rootCommand.AddOption(dryRunOption)

    let handler (dryRun: bool) =
        printfn $"%b{dryRun}"

    rootCommand.SetHandler(handler, dryRunOption)

    rootCommand.Invoke(args)
