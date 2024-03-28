module RobocopyLogs

open System.IO
open System.Text.RegularExpressions

let lineRegex = new Regex("\s*Files :\s+\d+\s+(\d+)\s+")

//let line = "    Files :      5117         0      5117         0         0         0"
//let line = "    Files :      5117         123      5117         0         0         0"
let fileHasCopies (fileName: string) =
    let mutable foundCopiesLine = false

    for line in File.ReadAllLines fileName do
        if not foundCopiesLine then
            let matches = lineRegex.Matches(line)

            if matches.Count > 0 then
                if "0" <> matches[0].Groups[1].Value then
                    foundCopiesLine <- true

    foundCopiesLine