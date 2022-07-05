module Tests

open Xunit

open ResetPerms.Cygwin

[<Fact>]
let ``Can convert cygwin path`` () =
    let input = "/cygdrive/c/Users/rober/OneDrive/data/Nightly Cron Job Times.xlsx"
    let output =  convertCygwinToWindows input
    
    Assert.True(output.IsSome)
    Assert.Equal(@"c:\Users\rober\OneDrive\data\Nightly Cron Job Times.xlsx", output.Value)

let changeDirLine = "rsync: [sender] change_dir \"/cygdrive/z/backup\" failed: Permission denied (13)"

[<Fact>]
let ``Can parse log line`` () =
    let output = parseLogLine changeDirLine

    Assert.True(output.IsSome)
    let rsyncLogLine = output.Value
    Assert.Equal("sender", rsyncLogLine.Actor)
    Assert.Equal("change_dir", rsyncLogLine.Message)
    Assert.Equal("/cygdrive/z/backup", rsyncLogLine.FileName)
    Assert.Equal("Permission denied (13)", rsyncLogLine.ErrorMessage)

[<Fact>]
let ``Can extract file name from log line`` () =
    let output = extractFailedFile changeDirLine
    
    Assert.True(output.IsSome)
    Assert.Equal("/cygdrive/z/backup", output.Value)


