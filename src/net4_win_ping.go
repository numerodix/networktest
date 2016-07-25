package main

import "errors"
import "fmt"
import "regexp"
import "strings"
import "strconv"


type WinPinger4 struct {
    ctx AppContext
}


func NewWinPinger4(ctx AppContext) WinPinger4 {
    return WinPinger4{
        ctx: ctx,
    }
}


func (wi WinPinger4) ping(host string, cnt int, timeoutMs int) PingExecution {
    var mgr = ProcMgr("ping",
                        "-n", fmt.Sprintf("%d", cnt),
                        "-w", fmt.Sprintf("%d", timeoutMs / 1000),
                        host)

    // DISABLE: Seems to kill the process prematurely on Windows
    //mgr.timeoutMs = timeoutMs

    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        wi.ctx.ft.printError("Failed to invoke ping", res.err)
        return PingExecution{Err: res.err}
    }

    // Extract the output
    var pingExec = wi.parsePing4(res.stdout)

    // Parsing failed :(
    if pingExec.Err != nil {
        wi.ctx.ft.printError("Failed to parse ipv4 network info", pingExec.Err)
    }

    return pingExec
}


func (wi WinPinger4) parsePing4(stdout string) PingExecution {
    /* Output:
      C:\> ping -n 1 -w 2 yahoo.com
      
      Pinging yahoo.com [206.190.36.45] with 32 bytes of data:
      Reply from 206.190.36.45: bytes=32 time=193ms TTL=51
      
      Ping statistics for 206.190.36.45:
          Packets: Sent = 1, Received = 1, Lost = 0 (0% loss),
      Approximate round trip times in milli-seconds:
          Minimum = 193ms, Maximum = 193ms, Average = 193ms
    */

    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxHost := regexp.MustCompile("^Pinging ([^ ]+)")
    rxStats := regexp.MustCompile("Average = ([0-9.]+)ms")

    // Loop variables
    var host = ""
    var timeAvg = -1.0
    var err error

    for _, line := range lines {
        if rxHost.MatchString(line) {
            host = rxHost.FindStringSubmatch(line)[1]
        }

        if rxStats.MatchString(line) {
            var timeA = rxStats.FindStringSubmatch(line)[1]
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
