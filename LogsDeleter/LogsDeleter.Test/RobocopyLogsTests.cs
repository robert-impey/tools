using LogsDeleter.RemoveEmptyRobocopyLogs;
using Shouldly;

namespace LogsDeleter.Test;

public class RobocopyLogsTests
{
    // The class containing the logic (RobocopyLogs) must be available.
    // Assuming the function is a static method:
    // public static bool IsFilesCopiedLine(string line) { ... }

    [Fact]
    public void NoFilesCopied()
    {
        var line = "    Files :       5117          0       5117          0          0          0";

        RobocopyLogsParser.IsFilesCopiedLine(line).ShouldBeFalse();
    }

    [Fact]
    public void FilesCopied()
    {
        var line = "    Files :       5117          123       5117          0          0          0";

        RobocopyLogsParser.IsFilesCopiedLine(line).ShouldBeTrue();
    }
}
