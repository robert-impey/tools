using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using Spectre.Console.Cli;
using Tools.Lib;

Console.OutputEncoding = System.Text.Encoding.UTF8;

var services = new ServiceCollection();
services.AddLogging(configure =>
{
    configure.AddConsole();
    configure.SetMinimumLevel(LogLevel.Information);
});

var registrar = new ServiceCollectionRegistrar(services);

var app = new CommandApp(registrar);

app.Configure(config =>
{
    config.SetApplicationName("TidyFolder");

    config.AddCommand<TidyFolder.Search.Command>("search")
        .WithDescription("Searches a single directory");

});

await app.RunAsync(args);

