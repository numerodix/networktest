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


func Test_linuxParseIpAddr4(t *testing.T) {
    var info = IP4NetworkInfo{}

    linuxParseIpAddr4(ip4AddrOutput, &info)

    // Errors
    assertIntEq(t, 0, len(info.Errs), "Errs does not match")

    // Networks
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
