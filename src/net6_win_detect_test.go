package main

import "testing"


const ipconfig6Output = `
Ethernet adapter Local Area Connection:

   Connection-specific DNS Suffix  . : ecoast.example.com
   IPv6 Address. . . . . . . . . . . : 2001:db8:21da:7:713e:a426:d167:37ab
   Temporary IPv6 Address. . . . . . : 2001:db8:21da:7:5099:ba54:9881:2e54
   Link-local IPv6 Address . . . . . : fe80::713e:a426:d167:37ab%6
   IPv4 Address. . . . . . . . . . . : 157.60.14.11
   Subnet Mask . . . . . . . . . . . : 255.255.255.0
   Default Gateway . . . . . . . . . : fe80::20a:42ff:feb0:5400%6
                                       157.60.14.1

Tunnel adapter Local Area Connection* 6:

   Connection-specific DNS Suffix  . : ecoast.example.com
   IPv6 Address. . . . . . . . . . . : 2001:db8:908c:f70f:200:5efe:157.60.14.11
   Link-local IPv6 Address . . . . . : fe80::200:5efe:157.60.14.11%9
   Default Gateway . . . . . . . . . : fe80::200:5efe:131.107.25.1%9

Tunnel adapter Local Area Connection* 7:

   Media State . . . . . . . . . . . : Media disconnected
   Connection-specific DNS Suffix  . :
`


func Test_winParseIpconfig6(t *testing.T) {
    var info = IPNetworkInfo{}

    var ctx = TestAppContext()
    var detector = NewWinNetDetect6(ctx)
    detector.parseIpconfig6(ipconfig6Output, &info)
    info.normalize()
    return /// XXX

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Nets
    assertIntEq(t, 2, len(info.Nets), "wrong number of gateways")

    assertStrEq(t, "eth1", info.Nets[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.0", info.Nets[0].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffff00", info.Nets[0].Ip.Mask.String(), "Mask does not match")

    assertStrEq(t, "wlan1", info.Nets[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.0", info.Nets[1].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffff00", info.Nets[1].Ip.Mask.String(), "Mask does not match")

    // Ips
    assertIntEq(t, 2, len(info.Ips), "wrong number of ips")

    assertStrEq(t, "eth1", info.Ips[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.11", info.Ips[0].Ip.String(), "Ip does not match")
    assertStrEq(t, "255.255.255.0", info.Ips[0].Mask.String(), "Mask does not match")

    assertStrEq(t, "wlan1", info.Ips[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.7", info.Ips[1].Ip.String(), "Ip does not match")
    assertStrEq(t, "255.255.255.0", info.Ips[1].Mask.String(), "Mask does not match")

    // Gws
    assertIntEq(t, 1, len(info.Gws), "wrong number of gateways")

    assertStrEq(t, "eth1", info.Gws[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.1", info.Gws[0].Ip.String(), "Ip does not match")

    // Ns hosts
    assertIntEq(t, 2, len(info.NsHosts), "wrong number of dns servers")

    assertStrEq(t, "192.168.1.1", info.NsHosts[0].Ip.String(), "Ip does not match")
    assertStrEq(t, "192.168.1.2", info.NsHosts[1].Ip.String(), "Ip does not match")
}

