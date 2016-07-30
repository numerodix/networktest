package main

import "net"
import "regexp"
import "strings"


type WinNetDetect4 struct {
    ctx AppContext
}


func NewWinNetDetect4(ctx AppContext) WinNetDetect4 {
    return WinNetDetect4{
        ctx: ctx,
    }
}


func (wnd WinNetDetect4) detectNetConn4() IPNetworkInfo {
    var info = IPNetworkInfo{}

    wnd.detectIpconfig4(&info)

    return info
}


func (wnd WinNetDetect4) detectIpconfig4(info *IPNetworkInfo) {
    var mgr = ProcMgr("ipconfig", "/all")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        wnd.ctx.ft.printError("Failed to detect ipv4 network", res.err)
        return
    }

    // Extract the output
    wnd.parseIpconfig4(res.stdout, info)

    // Parsing failed :(
    wnd.ctx.ft.printErrors("Failed to parse ipv4 network info", info.Errs)
}


func (wnd WinNetDetect4) parseIpconfig4(stdout string, info *IPNetworkInfo) {
    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxSection := regexp.MustCompile("^([^ ].+):")
    rxIP4Addr := regexp.MustCompile("[ ]{3}IPv4 Address.*: ([0-9.]+)")
    rxSubnet := regexp.MustCompile("[ ]{3}Subnet Mask.*: ([0-9.]+)")
    rxGw := regexp.MustCompile("[ ]{3}Default Gateway.*: ([0-9.]+)")
    rxDnsServer1 := regexp.MustCompile("[ ]{3}DNS Servers.*: ([0-9.]+)")
    rxDnsServern := regexp.MustCompile("[ ]{39}([0-9.]+)")

    // Loop variables
    var section = ""
    var ip = ""
    var subnet = ""
    var gw = ""
    var dnss []string

    var namer = InterfaceNamer()
    var sectionId = 0
    var inSection = false
    var inDns = false

    for _, line := range lines {
        if rxSection.MatchString(line) {
            // Terminate the previous section
            if inSection {
                if ip != "" {
                    var iface = namer.allocateName(section)

                    var ipobj = net.ParseIP(ip)
                    var maskobj = net.ParseIP(subnet)

                    var ipnet = ipAndMaskToIPNet(&ipobj, &maskobj)
                    var netMasked = applyMask(&ipnet.IP, &ipnet.Mask)

                    var gwobj = net.ParseIP(gw)

                    // Populate info
                    info.Nets = append(info.Nets, Network{
                        Iface: Interface{Name: iface},
                        Ip: netMasked,
                    })
                    info.Ips = append(info.Ips, IpAddr{
                        Iface: Interface{Name: iface},
                        Ip: ipobj,
                        Mask: maskobj,
                    })
                    info.Gws = append(info.Gws, Gateway{
                        Iface: Interface{Name: iface},
                        Ip: gwobj,
                    })

                    for _, dns := range dnss {
                        var nsobj = net.ParseIP(dns)

                        info.NsHosts = append(info.NsHosts, NsServer{
                            Ip: nsobj,
                        })
                    }
                }

                // Reset loop variables
                section, ip, subnet, gw, inDns = "", "", "", "", false
            }

            section = rxSection.FindStringSubmatch(line)[1]
            inSection = true
            sectionId += 1
        }

        if inSection {
            if rxIP4Addr.MatchString(line) {
                ip = rxIP4Addr.FindStringSubmatch(line)[1]
                inDns = false
            }
            if rxSubnet.MatchString(line) {
                subnet = rxSubnet.FindStringSubmatch(line)[1]
                inDns = false
            }
            if rxGw.MatchString(line) {
                gw = rxGw.FindStringSubmatch(line)[1]
                inDns = false
            }

            if rxDnsServer1.MatchString(line) {
                var dns = rxDnsServer1.FindStringSubmatch(line)[1]
                dnss = append(dnss, dns)
                inDns = true
            }

            if inDns && rxDnsServern.MatchString(line) {
                var dns = rxDnsServern.FindStringSubmatch(line)[1]
                dnss = append(dnss, dns)
            }
        }
    }
}
