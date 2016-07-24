package main


func pingJob(pinger LinuxPinger4, ch chan PingExecution,
             host string, cnt int, timeoutMs int) {

    var pingExec = pinger.ping(host, cnt, timeoutMs)
    ch <- pingExec
}


func runPings(pinger LinuxPinger4,
              hosts []string, cnt int, timeoutMs int) Pings {

    var pings = make(map[string]PingExecution)
    ch := make(chan PingExecution)

    // Launch
    for _, host := range hosts {
        go pingJob(pinger, ch, host, cnt, timeoutMs)
    }

    // Collect
    for _, host := range hosts {
        var pingExec = <-ch
        pings[pingExec.Host] = pingExec

        if host == "" {}  // "use" host variable to get around compiler error
    }

    return pings
}
