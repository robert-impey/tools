using Spectre.Console.Cli;

var app = new CommandApp<GenerateWindowsConfigSynchScripts.DefaultCommand>();
await app.RunAsync(args);
