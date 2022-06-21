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