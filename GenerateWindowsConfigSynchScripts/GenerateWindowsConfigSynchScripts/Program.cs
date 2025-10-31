using FolderManager;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;

var services = new ServiceCollection();
services.AddLogging(configure =>
{
    configure.AddConsole();
    configure.SetMinimumLevel(LogLevel.Information);
});
var registrar = new ServiceCollectionRegistrar(services);

var app = new CommandApp<GenerateWindowsConfigSynchScripts.DefaultCommand>(registrar);
await app.RunAsync(args);
