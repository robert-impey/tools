module RobocopyLogs

open System.IO
open System.Text.RegularExpressions

let lineRegex = new Regex("\s*Files :\s+\d+\s+(\d+)\s+")

let isFilesCopiedLine (line: string) =
    let matches = lineRegex.Matches(line)
    
    matches.Count > 0 && "0" <> matches[0].Groups[1].Value 

let fileHasCopies (fileName: string) =
    File.ReadAllLines fileName |> Array.exists isFilesCopiedLine
