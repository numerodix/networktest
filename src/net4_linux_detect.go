package main

import "fmt"
import "net"
import "regexp"
import "strings"


type LinuxNetDetect4 struct {
    col ColorBrush
    ft Formatter
}


func LinuxNetworkDetector4(col ColorBrush, ft Formatter) LinuxNetDetect4 {
    return LinuxNetDetect4{
        col: col,
        ft: ft,
    }
}


func (lnd *LinuxNetDetect4) linuxDetectNetConn4() IP4NetworkInfo {
    var info = IP4NetworkInfo{}

    lnd.linuxDetectIpAddr4(&info)
    lnd.linuxDetectIpRoute4(&info)
    unixDetectNsHosts4(&info)

    return info
}


func (lnd *LinuxNetDetect4) linuxDetectIpAddr4(info *IP4NetworkInfo) {
    var mgr = ProcMgr("/sbin/ip", "-4", "addr", "show")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        lnd.ft.printError("Failed to detect ipv4 network", res.err)
        return
    }

    // Extract the output
    lnd.linuxParseIpAddr4(res.stdout, info)
}

func (lnd *LinuxNetDetect4) linuxDetectIpRoute4(info *IP4NetworkInfo) {
    var mgr = ProcMgr("/sbin/ip", "-4", "route", "show")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        lnd.ft.printError("Failed to detect ipv4 routes", res.err)
        return
    }

    // Extract the output
    lnd.linuxParseIpRoute4(res.stdout, info)
}


func (lnd *LinuxNetDetect4) linuxParseIpAddr4(stdout string, info *IP4NetworkInfo) {
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

    for _, line := range lines {
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

            var netMasked = ipMaskToNet4(&ipnet.IP, &ipnet.Mask)
            var mask = ipnetToMask4(ipnet)

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


func (lnd *LinuxNetDetect4) linuxParseIpRoute4(stdout string, info *IP4NetworkInfo) {
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
