package main

import "errors"
import "fmt"
import "regexp"
import "strings"
import "strconv"


type LinuxPinger4 struct {
    ft Formatter
}


func NewLinuxPinger4(ft Formatter) LinuxPinger4 {
    return LinuxPinger4{
        ft: ft,
    }
}


func (pi *LinuxPinger4) ping(host string, cnt int, timeoutMs int) PingExecution {
    var mgr = ProcMgr("ping",
                        host,
                        fmt.Sprintf("-c%d", cnt),
                        fmt.Sprintf("-W%d", timeoutMs))
    mgr.timeoutMs = timeoutMs
    var res = mgr.run()

    // Use stderr as signal since err may not be reliable
    if res.stderr != "" {
        res.err = errors.New(strings.TrimSpace(res.stderr))
    }

    // The command failed :(
    if res.err != nil {
        pi.ft.printError("Failed to invoke ping", res.err)
        return PingExecution{Err: res.err}
    }

    // Extract the output
    var pingExec = pi.parsePing4(res.stdout)

    // Parsing failed :(
    if pingExec.Err != nil {
        pi.ft.printError("Failed to parse ipv4 network info", pingExec.Err)
    }

    return pingExec
}


func (pi *LinuxPinger4) parsePing4(stdout string) PingExecution {
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
    rxHost := regexp.MustCompile("^PING ([^ ]+)")
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
