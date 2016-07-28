package main


type NetDetector4 interface {
    detectNetConn4() IPNetworkInfo
}


func getDetector4(ctx AppContext) NetDetector4 {
    var det NetDetector4

    switch ctx.osName {
    // Linux userland
    case "linux":
        det = NewLinuxNetDetect4(ctx)

    // BSD userland
    case "darwin":
        fallthrough
    case "dragonfly":
        fallthrough
    case "freebsd":
        fallthrough
    case "netbsd":
        fallthrough
    case "openbsd":
        det = NewBsdNetDetect4(ctx)

    // Windows userland
    case "windows":
        det = NewWinNetDetect4(ctx)
    }

    return det
}
