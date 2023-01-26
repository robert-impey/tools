namespace Tools;

public static class ExecutablesHelper
{
    public static string FindExecutable(string name)
    {
        var isLinux = OsHelper.IsLinux();
        var isWindows = OsHelper.IsWindows();

        // Where is the executable?
        string? executablesDir = null;

        if (isWindows)
        {
            executablesDir = Environment.GetEnvironmentVariable("EXECUTABLES");
        }

        if (isLinux)
        {
            executablesDir = Path.Join(Environment.GetEnvironmentVariable("HOME"), "executables", "Linux", "prod", "x64");
        }

        string? executable = null;

        var exeSearch = Path.Join(executablesDir, name);

        if (isWindows)
        {
            exeSearch = $"{exeSearch}.exe";
        }

        if (File.Exists(exeSearch))
        {
            executable = exeSearch;
        }
   
        if (executable is null)
        {
            throw new Exception($"{name} executable not found!");
        }

        return executable;
    }
}
