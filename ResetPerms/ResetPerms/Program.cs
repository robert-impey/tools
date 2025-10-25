using ResetPerms;
using Spectre.Console.Cli;
using System.Runtime.InteropServices;

if (RuntimeInformation.IsOSPlatform(OSPlatform.Windows))
{
    Console.Error.WriteLine("This program should only be run on Linux or macOs");
    return 1;
}

var app = new CommandApp<DefaultCommand>();

app.Configure(config =>
{
    config.SetApplicationName("Reset Perms");
});

return app.Run(args);
