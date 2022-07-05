﻿module ResetPerms.Cygwin

let convertCygwinToWindows (file : string) =
    let prefix = "/cygdrive/"
    if file.StartsWith(prefix) then
        let startIndex = prefix.Length
        let shortened = file.Substring(startIndex)

        if shortened.[1] = '/' then
            let path = shortened.Substring(0, 1) + ":" + shortened.Substring(1)
            Some (path.Replace('/', '\\'))
        else
            None
    else
        None

type RsyncLogLine ={
    Actor: string
    Message: string
    FileName: string
    ErrorMessage: string
}

let parseLogLine (logLine: string) : Option<RsyncLogLine> =
    let rsyncPrefix = "rsync: ["
    if logLine.StartsWith(rsyncPrefix) then
        let indexStartActor = rsyncPrefix.Length
        let logLineFromStartActor = logLine.Substring(indexStartActor)
        let indexEndActor = logLineFromStartActor.IndexOf(']')
        if indexEndActor > 0 then
            let actor = logLineFromStartActor.Substring(0, indexEndActor)

            let logLineFromEndActor = logLineFromStartActor.Substring(indexEndActor + 2)

            let indexEndMessage = logLineFromEndActor.IndexOf('"')

            if indexEndMessage > 1 then
                let message = logLineFromEndActor.Substring(0, indexEndMessage - 1)

                let logLineFromEndMessage = logLineFromEndActor.Substring(indexEndMessage + 1)
                let indexEndFileName = logLineFromEndMessage.IndexOf('"')

                if indexEndFileName > 0 then
                    let fileName = logLineFromEndMessage.Substring(0, indexEndFileName)

                    let logLineFromEndFileName = logLineFromEndMessage.Substring(indexEndFileName + 1)
                    let failedStr = " failed: "

                    if logLineFromEndFileName.StartsWith(failedStr) then
                        let indexStartErrorMessage = failedStr.Length
                        let errorMessage = logLineFromEndFileName.Substring(indexStartErrorMessage)

                        Some { 
                            Actor = actor 
                            Message = message
                            FileName = fileName
                            ErrorMessage = errorMessage
                        }
                    else
                        None
                else
                    None
            else
                None
            
        else
            None
    else
        None

let extractFailedFile (logLine: string) =
    let parts = parseLogLine logLine

    match parts with
    | Some { FileName = fn } -> Some fn
    | None -> None