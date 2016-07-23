package main

import "testing"


const ipconfigOutput = `
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
`


func Test_winParseIpconfig4(t *testing.T) {
    var info = IP4NetworkInfo{}

    var ft = Formatter{}
    var detector = WindowsNetworkDetector4(ft)
    detector.parseIpconfig4(ipconfigOutput, &info)

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Ips
//    assertIntEq(t, 1, len(info.Ips), "wrong number of ips")

//    assertStrEq(t, "192.168.1.7", info.Ips[0].Ip.String(), "Ip does not match")

    // Gws
    assertIntEq(t, 1, len(info.Gws), "wrong number of gateways")

    assertStrEq(t, "if3", info.Gws[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.1", info.Gws[0].Ip.String(), "Ip does not match")
}

