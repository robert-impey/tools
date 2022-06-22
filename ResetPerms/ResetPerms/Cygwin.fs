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

let extractFailedFile (file: string) =
    let senderPrefix =
        "rsync: [sender] send_files failed to open \""

    let senderPostfix =
        "\": Permission denied (13)"

    let generatorPrefix =
        "rsync: [generator] failed to set permissions on \""

    let generatorPostfix =
        ".\": Permission denied (13)"

    if
        file.StartsWith(senderPrefix)
        && file.EndsWith(senderPostfix)
    then
        let startIndex = senderPrefix.Length

        let endIndex =
            file.Length
            - (senderPrefix.Length + senderPostfix.Length)

        Some(file.Substring(startIndex, endIndex))
    elif
        file.StartsWith(generatorPrefix)
        && file.EndsWith(generatorPostfix)
    then
        let startIndex = generatorPrefix.Length

        let endIndex =
            file.Length
            - (generatorPrefix.Length + generatorPostfix.Length)

        Some(file.Substring(startIndex, endIndex))
    else
        None
