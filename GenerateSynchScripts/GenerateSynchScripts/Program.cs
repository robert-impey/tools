using GenerateSynchScripts.RoboCopy.Files;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;
using Tools.Lib;

var services = new ServiceCollection();
services.AddLogging(configure =>
{
    configure.AddConsole();
    configure.SetMinimumLevel(LogLevel.Information);
});

services.AddSingleton<ScriptGenerator>();

var registrar = new ServiceCollectionRegistrar(services);

var app = new CommandApp(registrar);

app.Configure(config =>
{
    config.SetApplicationName("GenerateSynchScripts");

    config.AddCommand<RcifCommand>("rcif")
        .WithDescription("Robocopy scripts for individual files");
});

await app.RunAsync(args);
