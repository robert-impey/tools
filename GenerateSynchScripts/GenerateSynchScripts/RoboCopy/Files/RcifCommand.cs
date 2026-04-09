using Spectre.Console.Cli;

namespace GenerateSynchScripts.RoboCopy.Files;

internal class RcifCommand : AsyncCommand<CommandSettings>
{
    private readonly ScriptGenerator _generator;

    public RcifCommand(ScriptGenerator generator)
    {
        ArgumentNullException.ThrowIfNull(generator);

        _generator = generator;
    }

    protected override async Task<int> ExecuteAsync(
        CommandContext context,
        CommandSettings settings,
        CancellationToken cancellationToken
        )
    {
        ArgumentNullException.ThrowIfNull(settings);

        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Script);
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Autogen);
        ArgumentException.ThrowIfNullOrWhiteSpace(settings.Files);

        var synchFile = await SynchFileParser.ParseFile(settings.Files);

        await _generator.Generate(
            synchFile.Id,
            settings.Autogen,
            settings.Script,
            synchFile.Source,
            synchFile.Destination,
            synchFile.Files);

        return 0;
    }
}
