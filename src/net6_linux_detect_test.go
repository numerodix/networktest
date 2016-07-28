package main

//import "fmt"
//import "strings"
import "testing"


const ip6AddrOutput = `
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 16436 
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qlen 1000
    inet6 2a00:a21f:41::a13/64 scope global 
       valid_lft forever preferred_lft forever
    inet6 fe80::a31:49ef:12ce:f5a1/64 scope link 
       valid_lft forever preferred_lft forever
`


const ip6RouteOutput = `
2a00:a21f:41::/64 dev eth0  proto kernel  metric 256 
fe80::/64 dev eth0  proto kernel  metric 256 
default via 2a00:a21f:41::1 dev eth0  metric 1024 
`


func Test_linuxParseIpAddr6(t *testing.T) {
    var info = IPNetworkInfo{}

    var ctx = TestAppContext()
    var detector = NewLinuxNetDetect6(ctx)
    detector.parseIpAddr6(ip6AddrOutput, &info)
    info.normalize()

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Networks
    assertIntEq(t, 3, len(info.Nets), "wrong number of networks")

    assertStrEq(t, "lo", info.Nets[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "::1", info.Nets[0].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffffffffffffffffffffffffffffff",
                    info.Nets[0].Ip.Mask.String(), "Mask does not match")

    assertStrEq(t, "eth0", info.Nets[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "2a00:a21f:41::", info.Nets[1].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffffffffffffff0000000000000000",
                    info.Nets[1].Ip.Mask.String(), "Mask does not match")

    assertStrEq(t, "eth0", info.Nets[2].Iface.Name, "Iface does not match")
    assertStrEq(t, "fe80::", info.Nets[2].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffffffffffffff0000000000000000",
                    info.Nets[2].Ip.Mask.String(), "Mask does not match")

    // Ips
    assertIntEq(t, 3, len(info.Ips), "wrong number of ips")

    assertStrEq(t, "lo", info.Ips[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "::1", info.Ips[0].Ip.String(), "Ip does not match")
    assertStrEq(t, "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
                    info.Ips[0].Mask.String(), "Mask does not match")

    assertStrEq(t, "eth0", info.Ips[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "2a00:a21f:41::a13", info.Ips[1].Ip.String(), "Ip does not match")
    assertStrEq(t, "ffff:ffff:ffff:ffff::", info.Ips[1].Mask.String(), "Mask does not match")

    assertStrEq(t, "eth0", info.Ips[2].Iface.Name, "Iface does not match")
    assertStrEq(t, "fe80::a31:49ef:12ce:f5a1", info.Ips[2].Ip.String(), "Ip does not match")
    assertStrEq(t, "ffff:ffff:ffff:ffff::", info.Ips[2].Mask.String(), "Mask does not match")
}


func Test_linuxParseIpRoute6(t *testing.T) {
    var info = IPNetworkInfo{}

    var ctx = TestAppContext()
    var detector = NewLinuxNetDetect6(ctx)
    detector.parseIpRoute6(ip6RouteOutput, &info)

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Gateways
    assertIntEq(t, 1, len(info.Gws), "wrong number of gws")

    assertStrEq(t, "eth0", info.Gws[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "2a00:a21f:41::1", info.Gws[0].Ip.String(), "Ip does not match")
}
