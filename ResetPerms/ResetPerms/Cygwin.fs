module ResetPerms.Cygwin

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
    Actor: string option
    Message: string option
    FileName: string option
    ErrorMessage: string option
}

let parseLogLine (logLine: string) : Option<RsyncLogLine> =
    let rsyncPrefix = "rsync: ["
    if logLine.StartsWith(rsyncPrefix) then
        let indexStartActor = rsyncPrefix.Length
        let logLineFromStartActor = logLine.Substring(indexStartActor)
        let indexEndActor = logLineFromStartActor.IndexOf(']')
        
        let noneRsyncLogLine = {
            Actor = None
            Message = None
            FileName = None
            ErrorMessage = None
        }
        
        if indexEndActor > 0 then
            let actor = logLineFromStartActor.Substring(0, indexEndActor)

            let logLineFromEndActor = logLineFromStartActor.Substring(indexEndActor + 2)

            let indexEndMessage = logLineFromEndActor.IndexOf('"')
            
            let actorRsyncLogLine = { noneRsyncLogLine with Actor = Some actor }
            
            if indexEndMessage > 1 then
                let message = logLineFromEndActor.Substring(0, indexEndMessage - 1)

                let messageRsyncLogLine = { actorRsyncLogLine with Message = Some message }
                
                let logLineFromEndMessage = logLineFromEndActor.Substring(indexEndMessage + 1)
                let indexEndFileName = logLineFromEndMessage.IndexOf('"')

                if indexEndFileName > 0 then
                    let fileName = logLineFromEndMessage.Substring(0, indexEndFileName)

                    let fileNameRsyncLogLine = { messageRsyncLogLine with FileName = Some fileName }
                    
                    let logLineFromEndFileName = logLineFromEndMessage.Substring(indexEndFileName + 1)
                    let indexColon = logLineFromEndFileName.IndexOf(':')

                    if indexColon >= 0 then
                        let errorMessage = logLineFromEndFileName.Substring(indexColon + 2)

                        Some { fileNameRsyncLogLine with ErrorMessage = Some errorMessage
                        }
                    else
                        Some fileNameRsyncLogLine
                else
                    Some messageRsyncLogLine
            else
                Some actorRsyncLogLine
        else
            None
    else
        None

let extractFailedFile (logLine: string) : string option =
    let parts = parseLogLine logLine

    match parts with
    | Some { FileName = fn } -> fn
    | None -> None