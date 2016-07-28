package main

//import "fmt"
//import "strings"
import "testing"


const ip4AddrOutput = `
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default 
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    inet 192.168.1.6/24 brd 192.168.1.255 scope global eth0
       valid_lft forever preferred_lft forever
3: wlan0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    inet 192.168.1.10/24 brd 192.168.1.255 scope global wlan0
       valid_lft forever preferred_lft forever
`


const ip4RouteOutput = `
default via 192.168.1.1 dev eth0  proto static 
192.168.1.0/24 dev eth0  proto kernel  scope link  src 192.168.1.6  metric 1 
192.168.1.0/24 dev wlan0  proto kernel  scope link  src 192.168.1.10  metric 9
`


func Test_linuxParseIpAddr4(t *testing.T) {
    var info = IPNetworkInfo{}

    var ctx = TestAppContext()
    var detector = NewLinuxNetDetect4(ctx)
    detector.parseIpAddr4(ip4AddrOutput, &info)
    info.normalize()

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Networks
    assertIntEq(t, 3, len(info.Nets), "wrong number of networks")

    assertStrEq(t, "lo", info.Nets[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "127.0.0.0", info.Nets[0].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ff000000", info.Nets[0].Ip.Mask.String(), "Mask does not match")

    assertStrEq(t, "eth0", info.Nets[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.0", info.Nets[1].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffff00", info.Nets[1].Ip.Mask.String(), "Mask does not match")

    assertStrEq(t, "wlan0", info.Nets[2].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.0", info.Nets[2].Ip.IP.String(), "Ip does not match")
    assertStrEq(t, "ffffff00", info.Nets[2].Ip.Mask.String(), "Mask does not match")

    // Ips
    assertIntEq(t, 3, len(info.Ips), "wrong number of ips")

    assertStrEq(t, "lo", info.Ips[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "127.0.0.1", info.Ips[0].Ip.String(), "Ip does not match")
    assertStrEq(t, "255.0.0.0", info.Ips[0].Mask.String(), "Mask does not match")

    assertStrEq(t, "eth0", info.Ips[1].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.6", info.Ips[1].Ip.String(), "Ip does not match")
    assertStrEq(t, "255.255.255.0", info.Ips[1].Mask.String(), "Mask does not match")

    assertStrEq(t, "wlan0", info.Ips[2].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.10", info.Ips[2].Ip.String(), "Ip does not match")
    assertStrEq(t, "255.255.255.0", info.Ips[2].Mask.String(), "Mask does not match")
}


func Test_linuxParseIpRoute4(t *testing.T) {
    var info = IPNetworkInfo{}

    var ctx = TestAppContext()
    var detector = NewLinuxNetDetect4(ctx)
    detector.parseIpRoute4(ip4RouteOutput, &info)

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Gateways
    assertIntEq(t, 1, len(info.Gws), "wrong number of gws")

    assertStrEq(t, "eth0", info.Gws[0].Iface.Name, "Iface does not match")
    assertStrEq(t, "192.168.1.1", info.Gws[0].Ip.String(), "Ip does not match")
}
