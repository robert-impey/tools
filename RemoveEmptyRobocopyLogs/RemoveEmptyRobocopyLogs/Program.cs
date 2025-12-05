using FolderManager;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using RemoveEmptyRobocopyLogs;
using Spectre.Console.Cli;

var services = new ServiceCollection();
services.AddLogging(configure =>
{
    configure.AddConsole();
    configure.SetMinimumLevel(LogLevel.Information);
});
var registrar = new ServiceCollectionRegistrar(services);

var app = new CommandApp<DefaultCommand>(registrar);

app.Configure(config =>
{
    config.SetApplicationName("RemoveEmptyRobocopyLogs");
    config.AddCommand<DefaultCommand>("default")
        .WithDescription("Removes Robocopy log files that report zero copied files.");
});

await app.RunAsync(args);
