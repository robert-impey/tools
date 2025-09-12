using Spectre.Console.Cli;

var app = new CommandApp<GenerateWindowsConfigSynchScripts.Command>();
await app.RunAsync(args);
