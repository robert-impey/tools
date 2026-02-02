using System.Runtime.InteropServices;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using ResetPerms;
using Spectre.Console.Cli;
using Tools.Lib;

if (RuntimeInformation.IsOSPlatform(OSPlatform.Windows))
{
    Console.Error.WriteLine("This program should only be run on Linux or macOs");
}
else
{
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

    await app.RunAsync(args);
}
