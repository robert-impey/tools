module Tests

open Xunit
open RobocopyLogs

[<Fact>]
let ``No files Copied`` () =
    let line = "    Files :      5117         0      5117         0         0         0"
    Assert.False(isFilesCopiedLine line)

let ``Files Copied`` () =
    let line = "    Files :      5117         123      5117         0         0         0"
    Assert.True(isFilesCopiedLine line)
    