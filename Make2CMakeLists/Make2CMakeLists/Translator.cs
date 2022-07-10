namespace Make2CMakeLists;

internal class Translator
{
    public static bool TranslateMakefile(FileInfo makefileFileInfo)
    {
        var makefile = MakefileParser.Parse(makefileFileInfo);

        if (makefile is null)
            return false;

        if (makefileFileInfo.Directory is null)
            return false;

        if (!Writers.CMakeLists.Write(makefile, makefileFileInfo.Directory))
            return false;

        if (!Writers.GitIgnore.Write(makefile, makefileFileInfo.Directory))
            return false;

        return true;
    }
}