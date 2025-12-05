using FolderManager;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using RunCliTests;
using Spectre.Console.Cli;

var services = new ServiceCollection();

services.AddLogging(configure =>
{
    configure.AddConsole();

    configure.SetMinimumLevel(LogLevel.Information);
});

services.AddTransient<DefaultCommand>();

var registrar = new ServiceCollectionRegistrar(services);

var app = new CommandApp<DefaultCommand>(registrar);

app.Configure(config =>
{
    config.SetApplicationName("TestRunnerCli");

    config.AddCommand<DefaultCommand>("default")
        .WithDescription("A custom test runner that finds and executes test files.");
});

await app.RunAsync(args);
