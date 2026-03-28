using LogsDeleter.RemoveEmptyRobocopyLogs;
using LogsDeleter.Sweeping;
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
services.AddSingleton<LogsDeleter.Sweeping.LogsDeleter>();

var registrar = new ServiceCollectionRegistrar(services);

var app = new CommandApp(registrar);

app.Configure(config =>
{
    config.SetApplicationName("LogsDeleter");

    config.AddCommand<RerlCommand>("rerl")
        .WithDescription("Removes Robocopy log files that report zero copied files.");

    config.AddCommand<SweepAllCommand>("sweepAll")
        .WithDescription("Sweep all the log directories");

    config.AddCommand<SweepFromCommand>("sweepFrom")
        .WithDescription("Sweep away the old log files for just one tool");
});

await app.RunAsync(args);
