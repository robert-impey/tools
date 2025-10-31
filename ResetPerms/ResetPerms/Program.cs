using FolderManager;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using ResetPerms;
using Spectre.Console.Cli;
using System.Runtime.InteropServices;

if (RuntimeInformation.IsOSPlatform(OSPlatform.Windows))
{
    Console.Error.WriteLine("This program should only be run on Linux or macOs");
    return 1;
}

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
    config.SetApplicationName("Reset Perms");
});

return app.Run(args);
