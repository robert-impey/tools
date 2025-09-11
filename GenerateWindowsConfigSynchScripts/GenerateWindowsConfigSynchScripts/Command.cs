using Spectre.Console.Cli;

namespace GenerateWindowsConfigSynchScripts;
public class Command : Command<CommandSettings>
{
    public override int Execute(CommandContext context, CommandSettings settings)
    {
        Console.WriteLine("Files: " + settings.Files);
        Console.WriteLine("Source: " + settings.Source);
        Console.WriteLine("Destination: " + settings.Destination);
        Console.WriteLine("Logs Directory: " + settings.LogsDirectory);
        Console.WriteLine("Autogen: " + settings.Autogen);
        return 0;
    }
}
