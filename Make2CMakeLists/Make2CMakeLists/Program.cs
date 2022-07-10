using System.Runtime.CompilerServices;

[assembly: InternalsVisibleTo("Make2CMakeLists.Tests")]

namespace Make2CMakeLists;

internal class Program
{
    public static void Main(string[] args)
    {
        foreach (var arg in args)
        {
            var makeFiles = Searcher.FindMakeFiles(new DirectoryInfo(arg));

            foreach (var makeFile in makeFiles)
            {
                if (Translator.TranslateMakefile(makeFile))
                {
                    Console.WriteLine($"Translated {makeFile.FullName}");
                }
                else
                {
                    Console.Error.WriteLine($"Could not translate {makeFile.FullName}!");
                }
            }
        }
    }
}