package main

//import "fmt"
import "net"
import "regexp"
import "strings"


type WinNetDetect6 struct {
    ctx AppContext
}


func NewWinNetDetect6(ctx AppContext) WinNetDetect6 {
    return WinNetDetect6{
        ctx: ctx,
    }
}


func (wnd WinNetDetect6) detectNetConn6() IPNetworkInfo {
    var info = IPNetworkInfo{}

    wnd.detectIpconfig6(&info)

    return info
}


func (wnd WinNetDetect6) detectIpconfig6(info *IPNetworkInfo) {
    var mgr = ProcMgr("ipconfig", "/all")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        wnd.ctx.ft.printError("Failed to detect ipv6 network", res.err)
        return
    }

    // Extract the output
    wnd.parseIpconfig6(res.stdout, info)

    // Parsing failed :(
    wnd.ctx.ft.printErrors("Failed to parse ipv6 network info", info.Errs)
}


func (wnd WinNetDetect6) parseIpconfig6(stdout string, info *IPNetworkInfo) {
    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxSection := regexp.MustCompile("^([^ ].+):")
    rxIP6Addr := regexp.MustCompile("[ ]{3}IPv6 Address.*: ([A-Fa-f0-9:.]+)")
    rxGw := regexp.MustCompile("[ ]{3}Default Gateway.*: ([A-Fa-f0-9:.]+)")
    rxDnsServer1 := regexp.MustCompile("[ ]{3}DNS Servers.*: ([A-Fa-f0-9:.]+)")
    rxDnsServern := regexp.MustCompile("[ ]{39}([A-Fa-f0-9:.]+)")

    // Loop variables
    var section = ""
    var ip = ""
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

                    var ipobj = ip6stringToIP(ip)
                    // we don't know what the subnet is :/
                    var maskobj = net.ParseIP("ffff:ffff:ffff:ffff::")

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

                    // A link-local ip is not a gateway
                    if !gwobj.IsLinkLocalUnicast() {
                        info.Gws = append(info.Gws, Gateway{
                            Iface: Interface{Name: iface},
                            Ip: gwobj,
                        })
                    }

                    for _, dns := range dnss {
                        var nsobj = net.ParseIP(dns)

                        info.NsHosts = append(info.NsHosts, NsServer{
                            Ip: nsobj,
                        })
                    }
                }

                // Reset loop variables
                section, ip, gw, inDns = "", "", "", false
            }

            section = rxSection.FindStringSubmatch(line)[1]
            inSection = true
            sectionId += 1
        }

        if inSection {
            if rxIP6Addr.MatchString(line) {
                ip = rxIP6Addr.FindStringSubmatch(line)[1]
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
