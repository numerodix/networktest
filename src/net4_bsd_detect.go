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


func (bnd *BsdNetDetect4) detectNetConn4() IP4NetworkInfo {
    var info = IP4NetworkInfo{}

    bnd.detectIfconfig4(&info)
    bnd.detectNetstat4(&info)

    var und = NewUnixNetDetect4(bnd.ctx)
    und.detectNsHosts4(&info)

    return info
}


func (bnd *BsdNetDetect4) detectIfconfig4(info *IP4NetworkInfo) {
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


func (bnd *BsdNetDetect4) detectNetstat4(info *IP4NetworkInfo) {
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
    bnd.ctx.ft.printErrors("Failed to parse ipv4 network info", info.Errs)
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

            var mask = maskBytesToMask4(maskBytes)
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


func (bnd *BsdNetDetect4) parseNetstat4(stdout string, info *IP4NetworkInfo) {
    /* Output:
      $ /sbin/netstat -n -r
      Routing tables
      
      Internet:
      Destination        Gateway            Flags      Netif Expire
      default            10.0.2.2           UGS         em0
      10.0.2.0/24        link#1             U           em0
      10.0.2.15          link#1             UHS         lo0
      127.0.0.1          link#2             UH          lo0
      
      Internet6:
      Destination                       Gateway                       Flags      Netif Expire
      ::/96                             ::1                           UGRS        lo0
      ::1                               link#2                        UH          lo0
      ::ffff:0.0.0.0/96                 ::1                           UGRS        lo0
      fe80::/10                         ::1                           UGRS        lo0
      fe80::%em0/64                     link#1                        U           em0
      fe80::a00:27ff:fef2:34a1%em0      link#1                        UHS         lo0
      fe80::%lo0/64                     link#2                        U           lo0
      fe80::1%lo0                       link#2                        UHS         lo0
      ff01::%em0/32                     fe80::a00:27ff:fef2:34a1%em0  U           em0
      ff01::%lo0/32                     ::1                           U           lo0
      ff02::/16                         ::1                           UGRS        lo0
      ff02::%em0/32                     fe80::a00:27ff:fef2:34a1%em0  U           em0
      ff02::%lo0/32                     ::1                           U           lo0
    */

    // We will read line by line
    var lines = strings.Split(stdout, "\n")

    // Prepare regex objects
    rxLabel4 := regexp.MustCompile("^Internet:")
    rxLabel6 := regexp.MustCompile("^Internet6:")
    rxFlags := regexp.MustCompile("^default[ \t]+([0-9.]+)[ \t]+([^ ]+)[ \t]+([^ ]+)")
    rxGw := regexp.MustCompile("G")

    // Loop variables
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

        if scope4 && rxFlags.MatchString(line) {
            ip = rxFlags.FindStringSubmatch(line)[1]
            flags = rxFlags.FindStringSubmatch(line)[2]
            iface = rxFlags.FindStringSubmatch(line)[3]

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
