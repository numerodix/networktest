package main

//import "fmt"
//import "strings"
import "testing"


const ip6RouteOutput = `
2a00:a21f:41::/64 dev eth0  proto kernel  metric 256 
fe80::/64 dev eth0  proto kernel  metric 256 
default via 2a00:a21f:41::1 dev eth0  metric 1024 
`


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
