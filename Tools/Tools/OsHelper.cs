using System.Runtime.InteropServices;

namespace Tools;

public static class OsHelper
{
    // See https://codepedia.info/dotnet-core-to-detect-operating-system-os-platform/
    public static bool IsLinux() => RuntimeInformation.IsOSPlatform(OSPlatform.Linux);
    public static bool IsWindows() => RuntimeInformation.IsOSPlatform(OSPlatform.Windows);

    public static string? GetHomeDir() => IsWindows() ? Environment.GetEnvironmentVariable("USERPROFILE") : Environment.GetEnvironmentVariable("HOME");
}
