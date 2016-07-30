package main

import "fmt"
import "net"
import "regexp"
import "strconv"
import "strings"


type BsdNetDetect6 struct {
    ctx AppContext
}


func NewBsdNetDetect6(ctx AppContext) BsdNetDetect6 {
    return BsdNetDetect6{
        ctx: ctx,
    }
}


func (bnd BsdNetDetect6) detectNetConn6() IPNetworkInfo {
    var info = IPNetworkInfo{}

    bnd.detectIfconfig6(&info)
    bnd.detectNetstat6(&info)

    var und = NewUnixNetDetect4(bnd.ctx)
    und.detectNsHosts4(&info)

    return info
}


func (bnd BsdNetDetect6) detectIfconfig6(info *IPNetworkInfo) {
    var mgr = ProcMgr("ifconfig")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        bnd.ctx.ft.printError("Failed to detect ipv6 network", res.err)
        return
    }

    // Extract the output
    bnd.parseIfconfig6(res.stdout, info)

    // Parsing failed :(
    bnd.ctx.ft.printErrors("Failed to parse ipv6 network info", info.Errs)
}


func (bnd BsdNetDetect6) detectNetstat6(info *IPNetworkInfo) {
    var mgr = ProcMgr("netstat", "-n", "-r")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        bnd.ctx.ft.printError("Failed to detect ipv6 routes", res.err)
        return
    }

    // Extract the output
    bnd.parseNetstat6(res.stdout, info)

    // Parsing failed :(
    bnd.ctx.ft.printErrors("Failed to parse ipv6 route info", info.Errs)
}


func (bnd BsdNetDetect6) parseIfconfig6(stdout string, info *IPNetworkInfo) {
    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxIface := regexp.MustCompile("^([^ ]+):")
    rxInet := regexp.MustCompile("^\tinet6 ([A-Fa-f0-9:]+).+prefixlen ([0-9]+)")

    // Loop variables
    var iface = ""
    var ip = ""
    var maskbits = ""

    for _, line := range lines {
        if rxIface.MatchString(line) {
            iface = rxIface.FindStringSubmatch(line)[1]
        }

        if rxInet.MatchString(line) {
            ip = rxInet.FindStringSubmatch(line)[1]
            maskbits = rxInet.FindStringSubmatch(line)[2]

            var maskBits, err = strconv.Atoi(maskbits)

            // Parse failed
            if err != nil {
                info.Errs = append(info.Errs, err)
                continue
            }

            var ipNet = fmt.Sprintf("%s/%d", ip, maskBits)
            var ipobj, ipnet, err2 = net.ParseCIDR(ipNet)

            // Parse failed
            if err2 != nil {
                info.Errs = append(info.Errs, err2)
                continue
            }

            var netMasked = applyMask(&ipnet.IP, &ipnet.Mask)
            var maskIp = ipnetMaskAsIP(ipnet)

            // Populate info
            info.Nets = append(info.Nets, Network{
                Iface: Interface{Name: iface},
                Ip: netMasked,
            })
            info.Ips = append(info.Ips, IpAddr{
                Iface: Interface{Name: iface},
                Ip: ipobj,
                Mask: maskIp,
            })
        }
    }
}


func (bnd BsdNetDetect6) parseNetstat6(stdout string, info *IPNetworkInfo) {
    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxLabel6 := regexp.MustCompile("^Internet6:")
    rxFlags := regexp.MustCompile("^default[ \t]+([A-Fa-f0-9:]+)[^ ]*[ \t]+([^ ]+)")
    rxNetif := regexp.MustCompile("Netif")
    rxGw := regexp.MustCompile("G")

    // Loop variables
    var netifOffset = -1
    var netifLength = len("Netif")
    var scope6 = false
    var iface = ""
    var ip = ""
    var flags = ""

    for _, line := range lines {
        if !scope6 && rxLabel6.MatchString(line) {
            scope6 = true
            continue
        }

        if scope6 && rxNetif.MatchString(line) {
            var beginEnd = rxNetif.FindStringIndex(line)
            netifOffset = beginEnd[0]
        }

        if scope6 && rxFlags.MatchString(line) {
            ip = rxFlags.FindStringSubmatch(line)[1]
            flags = rxFlags.FindStringSubmatch(line)[2]

            // Find the end offset for the Netif field
            var ifaceField string
            var endOffset = netifOffset + netifLength
            if endOffset > (len(line) - 1) {
                ifaceField = line[netifOffset:]
            } else {
                ifaceField = line[netifOffset:endOffset]
            }

            iface = strings.TrimSpace(ifaceField)

            if rxGw.MatchString(flags) {
                var ipobj = net.ParseIP(ip)

                // A link-local ip is not a gateway
                if !ipobj.IsLinkLocalUnicast() {
                    // Populate info
                    info.Gws = append(info.Gws, Gateway{
                        Iface: Interface{Name: iface},
                        Ip: ipobj,
                    })
                }
            }
        }
    }
}
