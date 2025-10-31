using Spectre.Console.Cli;

namespace GenerateWindowsConfigSynchScripts;

internal class DefaultCommand : AsyncCommand<CommandSettings>
{
    private readonly WindowsConfigScriptsGenerator _generator;

    public DefaultCommand(WindowsConfigScriptsGenerator generator)
    {
        ArgumentNullException.ThrowIfNull(generator);

        _generator = generator;
    }

    public override async Task<int> ExecuteAsync(CommandContext context, CommandSettings settings)
    {
        ArgumentNullException.ThrowIfNull(settings);

        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Script);
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Autogen);
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Files);

        var synchFile = await SynchFileParser.ParseFile(settings.Files);

        await _generator.Generate(
            autogen: settings.Autogen,
            script: settings.Script,
            source: synchFile.Source,
            destination: synchFile.Destination,
            files: synchFile.Files);

        return 0;
    }
}
