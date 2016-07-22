package main

import "encoding/hex"
import "fmt"
import "net"
import "regexp"
import "strings"


type BsdNetDetect4 struct {
    ft Formatter
}


func BsdNetworkDetector4(ft Formatter) BsdNetDetect4 {
    return BsdNetDetect4{
        ft: ft,
    }
}


func (bnd *BsdNetDetect4) detectNetConn4() IP4NetworkInfo {
    var info = IP4NetworkInfo{}

    bnd.detectIfconfig4(&info)
//    fnd.detectIpRoute4(&info)

    var und = UnixNetworkDetector4(bnd.ft)
    und.detectNsHosts4(&info)

    return info
}


func (bnd *BsdNetDetect4) detectIfconfig4(info *IP4NetworkInfo) {
    var mgr = ProcMgr("/sbin/ifconfig")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        bnd.ft.printError("Failed to detect ipv4 network", res.err)
        return
    }

    // Extract the output
    bnd.parseIfconfig4(res.stdout, info)

    // Parsing failed :(
    for _, err := range info.Errs {
        bnd.ft.printError("Failed to parse ipv4 network info", err)
        return
    }
}


func (bnd *BsdNetDetect4) parseIfconfig4(stdout string, info *IP4NetworkInfo) {
    /* Output:
      $ /sbin/ifconfig
      em0: flags=8843<UP,BROADCAST,RUNNING,SIMPLEX,MULTICAST> metric 0 mtu 1500
      	options=9b<RXCSUM,TXCSUM,VLAN_MTU,VLAN_HWTAGGING,VLAN_HWCSUM>
      	ether 08:00:27:f2:34:a1
      	inet6 fe80::a00:27ff:fef2:34a1%em0 prefixlen 64 scopeid 0x1 
      	inet 10.0.2.15 netmask 0xffffff00 broadcast 10.0.2.255 
      	nd6 options=23<PERFORMNUD,ACCEPT_RTADV,AUTO_LINKLOCAL>
      	media: Ethernet autoselect (1000baseT <full-duplex>)
      	status: active
      lo0: flags=8049<UP,LOOPBACK,RUNNING,MULTICAST> metric 0 mtu 16384
      	options=600003<RXCSUM,TXCSUM,RXCSUM_IPV6,TXCSUM_IPV6>
      	inet6 ::1 prefixlen 128 
      	inet6 fe80::1%lo0 prefixlen 64 scopeid 0x2 
      	inet 127.0.0.1 netmask 0xff000000 
      	nd6 options=21<PERFORMNUD,AUTO_LINKLOCAL>
    */

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

            var mask = maskBytesToMask(maskBytes)
            var maskBits, _ = mask.Size()

            var ipNet = fmt.Sprintf("%s/%d", ip, maskBits)
            var ipobj, ipnet, err2 = net.ParseCIDR(ipNet)

            // Parse failed
            if err2 != nil {
                info.Errs = append(info.Errs, err2)
                continue
            }

            var netMasked = ipMaskToNet4(&ipnet.IP, &ipnet.Mask)
            var maskIp = ipnetToMask4(ipnet)

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


