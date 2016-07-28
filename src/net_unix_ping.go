package main

import "errors"
import "fmt"
import "net"
import "regexp"
import "strings"
import "strconv"


type UnixPinger struct {
    ctx AppContext
}


func NewUnixPinger(ctx AppContext) UnixPinger {
    return UnixPinger{
        ctx: ctx,
    }
}


func (ui UnixPinger) getPingArgs(host string, cnt int,
                                  timeoutMs int) (bool, bool, string, []string) {

    var ip = net.ParseIP(host)
    var exe = "ping"  // default to ipv4 ping

    // If the ip is link local then it's not a valid ping target
    var pingable = !ip.IsLinkLocalUnicast()

    // If the ip is ipv6 use ipv6 ping
    if ip != nil && ipIs6(ip) {
        exe = "ping6"
    // Otherwise it's a hostname, so use the ipver mode we are in
    } else if ip == nil && ui.ctx.ipver == 6 {
        exe = "ping6"
    }

    var args []string

    // ping (-c -W) seems to be supported everywhere
    // ...but ping6 (-c -W) only on linux
    if exe == "ping" || (exe == "ping6" && ui.ctx.isLinuxUserland()) {
        args = []string{
            "-c", fmt.Sprintf("%d", cnt),
            "-W", fmt.Sprintf("%d", timeoutMs / 1000),
            host,
        }

    // ping6 on a bsd platform only supports -c
    } else if exe == "ping6" && ui.ctx.isBsdUserland() {
        args = []string{
            "-c", fmt.Sprintf("%d", cnt),
            host,
        }
    }

    // now that we know which one we'll use, check if it's there
    var havetool = ui.ctx.toolbox.haveTool(exe)

    return havetool, pingable, exe, args
}


func (ui UnixPinger) ping(host string, cnt int, timeoutMs int) PingExecution {
    // Build the argument string
    var havetool, pingable, exe, args = ui.getPingArgs(host, cnt, timeoutMs)

    // If the ping tool isn't available there is no point trying
    if !havetool{
        return PingExecution{
            Host: host,
        }
    }

    // If the host is not pingable there is no point in trying
    if !pingable {
        return PingExecution{
            Host: host,
            Unpingable: true,
        }
    }

    var mgr = ProcMgr(exe, args...)
    mgr.timeoutMs = timeoutMs
    var res = mgr.run()

    // Use stderr as signal since err may not be reliable
    if res.stderr != "" {
        res.err = errors.New(strings.TrimSpace(res.stderr))
    }

    // The command failed :(
    if res.err != nil {
        ui.ctx.ft.printError(fmt.Sprintf("Failed to invoke %s", exe), res.err)
        return PingExecution{Err: res.err}
    }

    // Extract the output
    var pingExec = ui.parsePing(res.stdout)

    // Parsing failed :(
    if pingExec.Err != nil {
        ui.ctx.ft.printError(fmt.Sprintf("Failed to parse %s info", exe), pingExec.Err)
    }

    return pingExec
}


func (ui UnixPinger) parsePing(stdout string) PingExecution {
    /* Output:
      $ ping -c1 -W2 yahoo.com
      PING yahoo.com (98.138.253.109) 56(84) bytes of data.
      64 bytes from ir1.fp.vip.ne1.yahoo.com (98.138.253.109): icmp_seq=1 ttl=48 time=154 ms
      
      --- yahoo.com ping statistics ---
      1 packets transmitted, 1 received, 0% packet loss, time 0ms
      rtt min/avg/max/mdev = 154.327/154.327/154.327/0.000 ms
    */

    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxHost := regexp.MustCompile("^PING ([^)( ]+)")
    rxHostBsd := regexp.MustCompile("^PING6.+ bytes[)] ([^ ]+)")
    rxStats := regexp.MustCompile("= ([0-9.]+)/([0-9.]+)/([0-9.]+)/([0-9.]+) ms")

    // Loop variables
    var host = ""
    var timeAvg = -1.0
    var err error

    for _, line := range lines {
        if rxHost.MatchString(line) {
            host = rxHost.FindStringSubmatch(line)[1]
        }

        if rxHostBsd.MatchString(line) {
            host = rxHostBsd.FindStringSubmatch(line)[1]
        }

        if rxStats.MatchString(line) {
            var timeA = rxStats.FindStringSubmatch(line)[2]
            timeAvg, err = strconv.ParseFloat(timeA, 64)
        }
    }

    var res = PingExecution{
        Host: host,
        Time: timeAvg,
        Err: err,
    }

    return res
}
