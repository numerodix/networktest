package main

//import "fmt"
import "net"
import "regexp"
import "strings"


type LinuxNetDetect6 struct {
    ctx AppContext
}


func NewLinuxNetDetect6(ctx AppContext) LinuxNetDetect6 {
    return LinuxNetDetect6{
        ctx: ctx,
    }
}


func (lnd LinuxNetDetect6) detectNetConn6() IPNetworkInfo {
    var info = IPNetworkInfo{}

//    lnd.detectIpAddr6(&info)
    lnd.detectIpRoute6(&info)

//    var und = NewUnixNetDetect4(lnd.ctx)
//    und.detectNsHosts4(&info)

    return info
}


func (lnd LinuxNetDetect6) detectIpRoute6(info *IPNetworkInfo) {
    var mgr = ProcMgr("ip", "-6", "route", "show")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        lnd.ctx.ft.printError("Failed to detect ipv6 routes", res.err)
        return
    }

    // Extract the output
    lnd.parseIpRoute6(res.stdout, info)

    // Parsing failed :(
    lnd.ctx.ft.printErrors("Failed to parse ipv6 route info", info.Errs)
}


func (lnd LinuxNetDetect6) parseIpRoute6(stdout string, info *IPNetworkInfo) {
    /* Output:
      $ /sbin/ip -6 route show
      2a00:a21f:41::/64 dev eth0  proto kernel  metric 256 
      fe80::/64 dev eth0  proto kernel  metric 256 
      default via 2a00:a21f:41::1 dev eth0  metric 1024 
    */

    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxIface := regexp.MustCompile("^default via ([A-Fa-f0-9:]+) dev ([^ ]+)")

    // loop variables
    var iface = ""
    var ip = ""

    for _, line := range lines {
        if rxIface.MatchString(line) {
            ip = rxIface.FindStringSubmatch(line)[1]
            iface = rxIface.FindStringSubmatch(line)[2]

            var ipobj = net.ParseIP(ip)

            // Populate info
            info.Gws = append(info.Gws, Gateway{
                Iface: Interface{Name: iface},
                Ip: ipobj,
            })
        }
    }
}
