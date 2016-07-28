package main

//import "fmt"
//import "strings"
import "testing"


const ifconfig6Output = `
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
`


const netstat6Output = `
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
`


func Test_bsdParseIfconfig6(t *testing.T) {
    var info = IPNetworkInfo{}

    var ctx = TestAppContext()
    var detector = NewBsdNetDetect6(ctx)
    detector.parseIfconfig6(ifconfig6Output, &info)
    info.normalize()

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Networks
    assertIntEq(t, 3, len(info.Nets), "wrong number of networks")

    assertStrEq(t, "em0", info.Nets[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "fe80::", info.Nets[0].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffffffffffffff0000000000000000",
                    info.Nets[0].Ip.Mask.String(), "Mask does not match")

    assertStrEq(t, "lo0", info.Nets[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "::1", info.Nets[1].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffffffffffffffffffffffffffffff",
                    info.Nets[1].Ip.Mask.String(), "Mask does not match")

    assertStrEq(t, "lo0", info.Nets[2].Iface.Name, "Iface does not match")
    assertStrEq(t, "fe80::", info.Nets[2].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffffffffffffff0000000000000000",
                    info.Nets[2].Ip.Mask.String(), "Mask does not match")

    // Ips
    assertIntEq(t, 3, len(info.Ips), "wrong number of ips")

    assertStrEq(t, "em0", info.Ips[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "fe80::a00:27ff:fef2:34a1", info.Ips[0].Ip.String(), "Ip does not match")
    assertStrEq(t, "ffff:ffff:ffff:ffff::", info.Ips[0].Mask.String(), "Mask does not match")

    assertStrEq(t, "lo0", info.Ips[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "::1", info.Ips[1].Ip.String(), "Ip does not match")
    assertStrEq(t, "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
                    info.Ips[1].Mask.String(), "Mask does not match")

    assertStrEq(t, "lo0", info.Ips[2].Iface.Name, "Iface does not match")
    assertStrEq(t, "fe80::1", info.Ips[2].Ip.String(), "Ip does not match")
    assertStrEq(t, "ffff:ffff:ffff:ffff::", info.Ips[2].Mask.String(), "Mask does not match")
}


func Test_bsdParseNetstat6(t *testing.T) {
    var info = IPNetworkInfo{}

    var ctx = TestAppContext()
    var detector = NewBsdNetDetect6(ctx)
    detector.parseNetstat6(netstat6Output, &info)

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Gateways
    assertIntEq(t, 0, len(info.Gws), "wrong number of gws")

//    assertStrEq(t, "em0", info.Gws[0].Iface.Name, "Iface does not match")
//    assertStrEq(t, "10.0.2.2", info.Gws[0].Ip.String(), "Ip does not match")
}

