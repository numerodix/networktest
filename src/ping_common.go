package main


type Pinger interface {
    ping(host string, cnt int, timeoutMs int) PingExecution
}


func getPinger(ctx AppContext) Pinger {
    var pinger Pinger

    switch ctx.osName {
    // Windows is a special case
    case "windows":
        pinger = NewWinPinger4(ctx)

    default:
        pinger = NewLinuxPinger4(ctx)
    }

    return pinger
}
