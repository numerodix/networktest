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


func (pi LinuxPinger) getPingArgs(host string, cnt int, timeoutMs int) (string, []string) {

    var ip = net.ParseIP(host)
    var exe = "ping"  // default to ipv4 ping

    // If the ip is ipv6 use ipv6 ping
    if ip != nil && ipIs6(ip) {
        exe = "ping6"
    // Otherwise it's a hostname, so use the ipver mode we are in
    } else if ip == nil && pi.ctx.ipver == 6 {
        exe = "ping6"
    }

    var args []string

    // ping (-c -W) seems to be supported everywhere
    // ...but ping6 (-c -W) only on linux
    if exe == "ping" || (exe == "ping6" && pi.ctx.isLinuxUserland()) {
        args = []string{
            "-c", fmt.Sprintf("%d", cnt),
            "-W", fmt.Sprintf("%d", timeoutMs / 1000),
            host,
        }

    // ping6 on a bsd platform only supports -c
    } else if exe == "ping6" && pi.ctx.isBsdUserland() {
        args = []string{
            "-c", fmt.Sprintf("%d", cnt),
            host,
        }
    }

    return exe, args
}


func (pi LinuxPinger) ping(host string, cnt int, timeoutMs int) PingExecution {
    // Build the argument string
    var exe, args = pi.getPingArgs(host, cnt, timeoutMs)

    var mgr = ProcMgr(exe, args...)
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
