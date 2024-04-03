using FluentAssertions;

namespace Make2CMakeLists.Tests;

[TestClass]
public class MakefileParserTests
{
    [TestMethod]
    public void CanParseMakefile()
    {
        var makefileFileInfo = new FileInfo(Path.Join("Makefiles", "Makefile"));

        var makefile = MakefileParser.Parse(makefileFileInfo);

        makefile.Should().NotBeNull();

#pragma warning disable CS8602 // Dereference of a possibly null reference.
        makefile.CppStandard.Should().Be(11);
#pragma warning restore CS8602 // Dereference of a possibly null reference.
        makefile.Src.Should().Be("thread_process_demo.cpp");
        makefile.Target.Should().Be("thread_process_demo");
    }

    [TestMethod]
    public void CanParseCxxFlagsLine()
    {
        var standard = MakefileParser.ParseCxxFlagsLine("CXXFLAGS = -Wall -std=c++11 -pthread");

        standard.Should().Be(11);
    }

    [TestMethod]
    public void CanParseVarLine()
    {
        MakefileParser
            .ParseVarLine("CXXFLAGS", "CXXFLAGS = -Wall - std = c++11 - pthread")
            .Should().Be("-Wall - std = c++11 - pthread");
        MakefileParser
            .ParseVarLine("TARGET", "TARGET	 = thread_process_demo")
            .Should().Be("thread_process_demo");
        MakefileParser
            .ParseVarLine("SRC", "SRC		 = thread_process_demo.cpp")
            .Should().Be("thread_process_demo.cpp");
    }
}