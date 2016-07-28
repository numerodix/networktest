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

    // The Linux pinger works on all Unix systems
    default:
        pinger = NewUnixPinger(ctx)
    }

    return pinger
}
