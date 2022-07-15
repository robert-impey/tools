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
let failedToOpenLine = "rsync: [sender] send_files failed to open \"/cygdrive/c/Users/rober/OneDrive/config/_Common/Notepad++/config.xml\": Permission denied (13)"

[<Fact>]
let ``Can parse log line`` () =
    let changeDirLineOutput = parseLogLine changeDirLine

    Assert.True(changeDirLineOutput.IsSome)
    let changeDirLineRsyncLogLine = changeDirLineOutput.Value
    Assert.Equal("sender", changeDirLineRsyncLogLine.Actor.Value)
    Assert.Equal("change_dir", changeDirLineRsyncLogLine.Message.Value)
    Assert.Equal("/cygdrive/z/backup", changeDirLineRsyncLogLine.FileName.Value)
    Assert.Equal("Permission denied (13)", changeDirLineRsyncLogLine.ErrorMessage.Value)
    
    let failedToOpenLineOutput = parseLogLine failedToOpenLine

    Assert.True(failedToOpenLineOutput.IsSome)
    let failedToOpenLineRsyncLogLine = failedToOpenLineOutput.Value
    Assert.Equal("sender", failedToOpenLineRsyncLogLine.Actor.Value)
    Assert.Equal("send_files failed to open", failedToOpenLineRsyncLogLine.Message.Value)
    Assert.Equal("/cygdrive/c/Users/rober/OneDrive/config/_Common/Notepad++/config.xml", failedToOpenLineRsyncLogLine.FileName.Value)
    Assert.Equal("Permission denied (13)", failedToOpenLineRsyncLogLine.ErrorMessage.Value)

[<Fact>]
let ``Can extract file name from log line`` () =
    let changeDirLineOutput = extractFailedFile changeDirLine
    
    Assert.True(changeDirLineOutput.IsSome)
    Assert.Equal("/cygdrive/z/backup", changeDirLineOutput.Value)

    let failedToOpenLineOutput = extractFailedFile failedToOpenLine
    
    Assert.True(failedToOpenLineOutput.IsSome)
    Assert.Equal("/cygdrive/c/Users/rober/OneDrive/config/_Common/Notepad++/config.xml", failedToOpenLineOutput.Value)

