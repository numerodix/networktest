package main

import "errors"
import "fmt"
import "net"
import "regexp"
import "strings"
import "strconv"


type LinuxPinger struct {
    ctx AppContext
}


func NewLinuxPinger(ctx AppContext) LinuxPinger {
    return LinuxPinger{
        ctx: ctx,
    }
}


func (pi LinuxPinger) getPingExecutable(host string) string {
    var ip = net.ParseIP(host)
    var exe = "ping"  // default to ipv4 ping

    // If the ip is ipv6 use ipv6 ping
    if ip != nil && ipIs6(ip) {
        exe = "ping6"
    // Otherwise it's a hostname, so use the ipver mode we are in
    } else if pi.ctx.ipver == 6 {
        exe = "ping6"
    }

    return exe
}


func (pi LinuxPinger) ping(host string, cnt int, timeoutMs int) PingExecution {
    // Decide whether to use ping or ping6
    var exe = pi.getPingExecutable(host)

    var mgr = ProcMgr(exe,
                        "-c", fmt.Sprintf("%d", cnt),
                        "-W", fmt.Sprintf("%d", timeoutMs / 1000),
                        host)
    mgr.timeoutMs = timeoutMs
    var res = mgr.run()

    // Use stderr as signal since err may not be reliable
    if res.stderr != "" {
        res.err = errors.New(strings.TrimSpace(res.stderr))
    }

    // The command failed :(
    if res.err != nil {
        pi.ctx.ft.printError(fmt.Sprintf("Failed to invoke %s", exe), res.err)
        return PingExecution{Err: res.err}
    }

    // Extract the output
    var pingExec = pi.parsePing(res.stdout)

    // Parsing failed :(
    if pingExec.Err != nil {
        pi.ctx.ft.printError(fmt.Sprintf("Failed to parse %s info", exe), pingExec.Err)
    }

    return pingExec
}


func (pi LinuxPinger) parsePing(stdout string) PingExecution {
    /* Output:
      $ ping -c1 -W2 yahoo.com
      PING yahoo.com (98.138.253.109) 56(84) bytes of data.
      64 bytes from ir1.fp.vip.ne1.yahoo.com (98.138.253.109): icmp_seq=1 ttl=48 time=154 ms
      
      --- yahoo.com ping statistics ---
      1 packets transmitted, 1 received, 0% packet loss, time 0ms
      rtt min/avg/max/mdev = 154.327/154.327/154.327/0.000 ms

      $ ping6 -c1 -W2 ::1
      PING ::1(::1) 56 data bytes
      64 bytes from ::1: icmp_seq=1 ttl=64 time=0.049 ms
      
      --- ::1 ping statistics ---
      1 packets transmitted, 1 received, 0% packet loss, time 0ms
      rtt min/avg/max/mdev = 0.049/0.049/0.049/0.000 ms
    */

    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxHost := regexp.MustCompile("^PING ([^)( ]+)")
    rxStats := regexp.MustCompile("= ([0-9.]+)/([0-9.]+)/([0-9.]+)/([0-9.]+) ms")

    // Loop variables
    var host = ""
    var timeAvg = -1.0
    var err error

    for _, line := range lines {
        if rxHost.MatchString(line) {
            host = rxHost.FindStringSubmatch(line)[1]
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
