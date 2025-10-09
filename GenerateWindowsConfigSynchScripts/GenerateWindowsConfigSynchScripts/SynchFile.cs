namespace GenerateWindowsConfigSynchScripts;

public record SynchFile(string Source, string Destination, IEnumerable<string> Files)
{
}
