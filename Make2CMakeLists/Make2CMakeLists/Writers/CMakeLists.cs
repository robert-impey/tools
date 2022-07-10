namespace Make2CMakeLists.Writers;

internal static class CMakeLists
{
    public static bool Write(Makefile makefile, DirectoryInfo directoryInfo)
    {
        var cMakeLists = Path.Join(directoryInfo.FullName, "CMakeLists.txt");

        using var outFile = File.CreateText(cMakeLists);

        outFile.WriteLine("cmake_minimum_required(VERSION 3.22)");
        outFile.WriteLine($"project({makefile.Target})");
        outFile.WriteLine();

        outFile.WriteLine($"set(CMAKE_CXX_STANDARD {makefile.CppStandard})");
        outFile.WriteLine();

        outFile.WriteLine($"add_executable({makefile.Target} {makefile.Src})");

        return true;
    }
}