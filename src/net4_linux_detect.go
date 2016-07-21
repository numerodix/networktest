package main

import "fmt"
import "net"
import "regexp"
import "strings"


func linuxDetectNetConn4() IP4NetworkInfo {
    var info = IP4NetworkInfo{}

    linuxDetectIpAddr4(&info)
    linuxDetectIpRoute4(&info)
    unixDetectNsHosts4(&info)

    return info
}


func linuxDetectIpAddr4(info *IP4NetworkInfo) {
    var mgr = ProcMgr("/sbin/ip", "-4", "addr", "show")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        // XXX print some kind of useful error
        return
    }

    // Extract the output
    linuxParseIpAddr4(res.stdout, info)
}

func linuxDetectIpRoute4(info *IP4NetworkInfo) {
    var mgr = ProcMgr("/sbin/ip", "-4", "route", "show")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        // XXX print some kind of useful error
        return
    }

    // Extract the output
    linuxParseIpRoute4(res.stdout, info)
}


func linuxParseIpAddr4(stdout string, info *IP4NetworkInfo) {
    /* Output:
      $ /sbin/ip -4 addr show
      1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default 
          inet 127.0.0.1/8 scope host lo
             valid_lft forever preferred_lft forever
      2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
          inet 192.168.1.6/24 brd 192.168.1.255 scope global eth0
             valid_lft forever preferred_lft forever
    */

    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxIface := regexp.MustCompile("^[0-9]+: ([^ ]+):")
    rxInet := regexp.MustCompile("^[ ]{4}inet ([0-9.]+)[/]([0-9]+)")

    // Loop variables
    var iface = ""
    var ip = ""
    var maskBits = ""

    for i := range lines {
        var line = lines[i]

        if rxIface.MatchString(line) {
            iface = rxIface.FindStringSubmatch(line)[1]
        }

        if rxInet.MatchString(line) {
            ip = rxInet.FindStringSubmatch(line)[1]
            maskBits = rxInet.FindStringSubmatch(line)[2]

            var ipNet = fmt.Sprintf("%s/%s", ip, maskBits)
            var ipobj, ipnet, err = net.ParseCIDR(ipNet)

            // Parse failed
            if err != nil {
                info.Errs = append(info.Errs, err)
                continue
            }

            var netMasked = ipMaskToNet4(ipnet.IP, ipnet.Mask)

            // XXX move to convenience function?
            var mask = net.IPv4(
                ipnet.Mask[0],
                ipnet.Mask[1],
                ipnet.Mask[2],
                ipnet.Mask[3])

            // Populate info
            info.Nets = append(info.Nets, Network{
                Iface: Interface{Name: iface},
                Ip: netMasked,
            })
            info.Ips = append(info.Ips, IpAddr{
                Iface: Interface{Name: iface},
                Ip: ipobj,
                Mask: mask,
            })
        }
    }
}


func linuxParseIpRoute4(stdout string, info *IP4NetworkInfo) {
    /* Output:
      $ /sbin/ip -4 route show
      default via 192.168.1.1 dev eth0  proto static 
      192.168.1.0/24 dev eth0  proto kernel  scope link  src 192.168.1.6  metric 1 
    */

    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxIface := regexp.MustCompile("^default via ([0-9.]+) dev ([^ ]+)")

    // loop variables
    var iface = ""
    var ip = ""

    for i := range lines {
        var line = lines[i]

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
