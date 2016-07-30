package main

import "encoding/hex"
import "fmt"
import "net"
import "regexp"
import "strings"


type BsdNetDetect4 struct {
    ctx AppContext
}


func NewBsdNetDetect4(ctx AppContext) BsdNetDetect4 {
    return BsdNetDetect4{
        ctx: ctx,
    }
}


func (bnd BsdNetDetect4) detectNetConn4() IPNetworkInfo {
    var info = IPNetworkInfo{}

    bnd.detectIfconfig4(&info)
    bnd.detectNetstat4(&info)

    var und = NewUnixNetDetect4(bnd.ctx)
    und.detectNsHosts4(&info)

    return info
}


func (bnd BsdNetDetect4) detectIfconfig4(info *IPNetworkInfo) {
    var mgr = ProcMgr("ifconfig")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        bnd.ctx.ft.printError("Failed to detect ipv4 network", res.err)
        return
    }

    // Extract the output
    bnd.parseIfconfig4(res.stdout, info)

    // Parsing failed :(
    bnd.ctx.ft.printErrors("Failed to parse ipv4 network info", info.Errs)
}


func (bnd BsdNetDetect4) detectNetstat4(info *IPNetworkInfo) {
    var mgr = ProcMgr("netstat", "-n", "-r")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        bnd.ctx.ft.printError("Failed to detect ipv4 routes", res.err)
        return
    }

    // Extract the output
    bnd.parseNetstat4(res.stdout, info)

    // Parsing failed :(
    bnd.ctx.ft.printErrors("Failed to parse ipv4 route info", info.Errs)
}


func (bnd BsdNetDetect4) parseIfconfig4(stdout string, info *IPNetworkInfo) {
    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxIface := regexp.MustCompile("^([^ ]+):")
    rxInet := regexp.MustCompile("^\tinet ([0-9.]+) netmask 0x([a-f0-9]+)")

    // Loop variables
    var iface = ""
    var ip = ""
    var maskHex = ""

    for _, line := range lines {
        if rxIface.MatchString(line) {
            iface = rxIface.FindStringSubmatch(line)[1]
        }

        if rxInet.MatchString(line) {
            ip = rxInet.FindStringSubmatch(line)[1]
            maskHex = rxInet.FindStringSubmatch(line)[2]

            var maskBytes, err = hex.DecodeString(maskHex)

            // Parse failed
            if err != nil {
                info.Errs = append(info.Errs, err)
                continue
            }

            var mask = maskBytesToIPMask4(maskBytes)
            var maskBits, _ = mask.Size()

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


func (bnd BsdNetDetect4) parseNetstat4(stdout string, info *IPNetworkInfo) {
    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxLabel4 := regexp.MustCompile("^Internet:")
    rxLabel6 := regexp.MustCompile("^Internet6:")
    rxFlags := regexp.MustCompile("^default[ \t]+([0-9.]+)[ \t]+([^ ]+)")
    rxNetif := regexp.MustCompile("Netif")
    rxGw := regexp.MustCompile("G")

    // Loop variables
    var netifOffset = -1
    var netifLength = len("Netif")
    var scope4 = false
    var iface = ""
    var ip = ""
    var flags = ""

    for _, line := range lines {
        if !scope4 && rxLabel4.MatchString(line) {
            scope4 = true
            continue
        }

        if scope4 && rxLabel6.MatchString(line) {
            scope4 = false
            break
        }

        if scope4 && rxNetif.MatchString(line) {
            var beginEnd = rxNetif.FindStringIndex(line)
            netifOffset = beginEnd[0]
        }

        if scope4 && rxFlags.MatchString(line) {
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

                // Populate info
                info.Gws = append(info.Gws, Gateway{
                    Iface: Interface{Name: iface},
                    Ip: ipobj,
                })
            }
        }
    }
}
