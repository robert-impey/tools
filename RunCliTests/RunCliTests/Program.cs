using FolderManager;
using Microsoft.Extensions.DependencyInjection;
using RunCliTests;
using Spectre.Console.Cli;

// This is a modernization of https://github.com/robert-impey/run-cli-tests

var services = new ServiceCollection();

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
