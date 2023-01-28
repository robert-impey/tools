namespace Tools;

public static class WellKnownFoldersHelper
{
    public static string GetLogsDir()
    {
        return Path.Join(OsHelper.GetHomeDir(), "logs");
    }
}