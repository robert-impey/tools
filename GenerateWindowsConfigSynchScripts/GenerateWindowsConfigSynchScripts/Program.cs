using GenerateWindowsConfigSynchScripts;
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

services.AddSingleton<WindowsConfigScriptsGenerator>();

var registrar = new ServiceCollectionRegistrar(services);

var app = new CommandApp<DefaultCommand>(registrar);

app.Configure(config =>
{
    config.SetApplicationName("GenerateWindowsConfigSynchScripts");
});

await app.RunAsync(args);
