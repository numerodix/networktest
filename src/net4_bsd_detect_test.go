package main

//import "fmt"
//import "strings"
import "testing"


const ifconfig4Output = `
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


const netstat4Output = `
`


func Test_bsdParseIpAddr4(t *testing.T) {
    var info = IP4NetworkInfo{}

    var ft = Formatter{}
    var detector = BsdNetworkDetector4(ft)
    detector.parseIpAddr4(ifconfig4Output, &info)

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Networks
    assertIntEq(t, 2, len(info.Nets), "wrong number of networks")

    assertStrEq(t, "em0", info.Nets[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "10.0.2.0", info.Nets[0].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffff00", info.Nets[0].Ip.Mask.String(), "Mask does not match")

    assertStrEq(t, "lo0", info.Nets[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "127.0.0.0", info.Nets[1].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ff000000", info.Nets[1].Ip.Mask.String(), "Mask does not match")

    // Ips
    assertIntEq(t, 2, len(info.Ips), "wrong number of ips")

    assertStrEq(t, "em0", info.Ips[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "10.0.2.15", info.Ips[0].Ip.String(), "Ip does not match")
    assertStrEq(t, "255.255.255.0", info.Ips[0].Mask.String(), "Mask does not match")

    assertStrEq(t, "lo0", info.Ips[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "127.0.0.1", info.Ips[1].Ip.String(), "Ip does not match")
    assertStrEq(t, "255.0.0.0", info.Ips[1].Mask.String(), "Mask does not match")
}
