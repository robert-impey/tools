namespace GenerateSynchScripts.RoboCopy.Files;

public record SynchFile(string Id, string Source, string Destination, IEnumerable<string> Files)
{
}
