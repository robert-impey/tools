module Tests

open Xunit

open ResetPerms.Cygwin

[<Fact>]
let ``My test`` () =
    let input = "/cygdrive/c/Users/rober/OneDrive/data/Nightly Cron Job Times.xlsx"
    let output =  convertCygwinToWindows input
    
    Assert.True(output.IsSome)
    Assert.Equal(@"c:\Users\rober\OneDrive\data\Nightly Cron Job Times.xlsx", output.Value)
    