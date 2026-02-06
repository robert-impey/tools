using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using RemoveEmptyRobocopyLogs;
using Spectre.Console.Cli;
using Tools.Lib;

var services = new ServiceCollection();
services.AddLogging(configure =>
{
    configure.AddConsole();
    configure.SetMinimumLevel(LogLevel.Information);
});
var registrar = new ServiceCollectionRegistrar(services);

// Use non-generic CommandApp so we can register commands from referenced libraries
var app = new CommandApp(registrar);

app.Configure(config =>
{
    config.SetApplicationName("LogsDeleter");

    // Register the RemoveEmptyRobocopyLogs subcommand(s)
    config.AddCommand<DefaultCommand>("rerl")
        .WithDescription("Removes Robocopy log files that report zero copied files.");

    // Additional subcommands can be registered here in future
});

await app.RunAsync(args);
