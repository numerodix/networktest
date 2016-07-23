package main

import "fmt"
import "net"
import "regexp"
import "strings"


type WinNetDetect4 struct {
    ft Formatter
}


func WindowsNetworkDetector4(ft Formatter) WinNetDetect4 {
    return WinNetDetect4{
        ft: ft,
    }
}


func (wnd *WinNetDetect4) detectNetConn4() IP4NetworkInfo {
    var info = IP4NetworkInfo{}

    wnd.detectIpconfig4(&info)

    return info
}


func (wnd *WinNetDetect4) detectIpconfig4(info *IP4NetworkInfo) {
    var mgr = ProcMgr("ipconfig")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        wnd.ft.printError("Failed to detect ipv4 network", res.err)
        return
    }

    // Extract the output
    wnd.parseIpconfig4(res.stdout, info)

    // Parsing failed :(
    wnd.ft.printErrors("Failed to parse ipv4 network info", info.Errs)
}


func (wnd *WinNetDetect4) parseIpconfig4(stdout string, info *IP4NetworkInfo) {
    /* Output:
      C:\> ipconfig /all
      Windows IP Configuration
       
         Host Name . . . . . . . . . . . . : DESKTOP-AB123XX
         Primary Dns Suffix  . . . . . . . :
         Node Type . . . . . . . . . . . . : Hybrid
         IP Routing Enabled. . . . . . . . : No
         WINS Proxy Enabled. . . . . . . . : No
       
      Ethernet adapter Ethernet:
       
         Media State . . . . . . . . . . . : Media disconnected
         Connection-specific DNS Suffix  . :
         Description . . . . . . . . . . . : Intel(R) 82579LM Gigabit Network Connection
         Physical Address. . . . . . . . . : FF-56-61-37-6C-31
         DHCP Enabled. . . . . . . . . . . : Yes
         Autoconfiguration Enabled . . . . : Yes
       
      Wireless LAN adapter Local Area Connection* 2:
       
         Media State . . . . . . . . . . . : Media disconnected
         Connection-specific DNS Suffix  . :
         Description . . . . . . . . . . . : Microsoft Wi-Fi Direct Virtual Adapter
         Physical Address. . . . . . . . . : FF-56-61-37-6C-31
         DHCP Enabled. . . . . . . . . . . : Yes
         Autoconfiguration Enabled . . . . : Yes
       
      Wireless LAN adapter Wi-Fi:
       
         Connection-specific DNS Suffix  . :
         Description . . . . . . . . . . . : Intel(R) Centrino(R) Advanced-N 6205
         Physical Address. . . . . . . . . : FF-56-61-37-6C-31
         DHCP Enabled. . . . . . . . . . . : Yes
         Autoconfiguration Enabled . . . . : Yes
         Link-local IPv6 Address . . . . . : fe80::bd72:d119:c03d:6033%11(Preferred)
         IPv4 Address. . . . . . . . . . . : 192.168.1.7(Preferred)
         Subnet Mask . . . . . . . . . . . : 255.255.255.0
         Lease Obtained. . . . . . . . . . : Saturday, September 06, 2009 12:28:03 AM
         Lease Expires . . . . . . . . . . : Saturday, September 06, 2009 12:28:03 AM
         Default Gateway . . . . . . . . . : 192.168.1.1
         DHCP Server . . . . . . . . . . . : 192.168.1.1
         DHCPv6 IAID . . . . . . . . . . . : 98606385
         DHCPv6 Client DUID. . . . . . . . : 00-01-00-01-12-31-BB-90-00-03-FF-16-46-11
         DNS Servers . . . . . . . . . . . : 192.168.1.1
                                             192.168.1.2
         NetBIOS over Tcpip. . . . . . . . : Enabled
       
      Ethernet adapter Bluetooth Network Connection:
       
         Media State . . . . . . . . . . . : Media disconnected
         Connection-specific DNS Suffix  . :
         Description . . . . . . . . . . . : Bluetooth Device (Personal Area Network)
         Physical Address. . . . . . . . . : FF-56-61-37-6C-31
         DHCP Enabled. . . . . . . . . . . : Yes
         Autoconfiguration Enabled . . . . : Yes
    */

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
    var ip = ""
    var subnet = ""
    var gw = ""
    var dnss []string

    var sectionId = 0
    var inSection = false
    var inDns = false

    for _, line := range lines {
        if rxSection.MatchString(line) {
            // Terminate the previous section
            if inSection {
                if ip != "" {
                    var iface = fmt.Sprintf("if%d", sectionId)

                    var ipobj = net.ParseIP(ip)
                    var maskobj = net.ParseIP(subnet)

                    var ipnet = ipIPMaskToNet4(&ipobj, &maskobj)
                    var netMasked = ipMaskToNet4(&ipnet.IP, &ipnet.Mask)

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
                ip, subnet, gw, inDns = "", "", "", false
            }

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
