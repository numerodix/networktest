package main

import (
    "fmt"
)


func DisplayLocalNetwork6(ft Formatter,
                          ip6Ips Ip6IpAddrExecution) {

    var ip6IpAddrBlocks = ip6Ips.Ip6IpAddrBlocks

    fmt.Printf("%s\n", ft.FormatHeader("Scanning for networks"))
    for i := range ip6IpAddrBlocks {
        var ip6IpAddrBlock = ip6IpAddrBlocks[i]

        var ifaceFmt = ft.FormatIfaceField(ip6IpAddrBlock.Iface)
        var netwFmt = ft.FormatIp6Field(ip6IpAddrBlock.Network.IP)
        var maskFmt = ft.FormatMask6Field(ip6IpAddrBlock.Network.Mask)
        var scopeFmt = ft.FormatScope6Field(ip6IpAddrBlock.Scope)
        fmt.Printf("    %s  %s %s   %s\n", ifaceFmt, netwFmt, maskFmt, scopeFmt)
    }
    if len(ip6IpAddrBlocks) == 0 {
        fmt.Printf("    %s\n", ft.FormatError("none found"))
    }

    fmt.Printf("%s\n", ft.FormatHeader("Detecting ips"))
    for i := range ip6IpAddrBlocks {
        var ip6IpAddrBlock = ip6IpAddrBlocks[i]

        var ifaceFmt = ft.FormatIfaceField(ip6IpAddrBlock.Iface)
        var ipFmt = ft.FormatIp6Field(ip6IpAddrBlock.IPv6)
        var maskFmt = ft.FormatMask6Field(ip6IpAddrBlock.Network.Mask)
        fmt.Printf("    %s  %s %s\n", ifaceFmt, ipFmt, maskFmt)
    }
    if len(ip6IpAddrBlocks) == 0 {
        fmt.Printf("    %s\n", ft.FormatError("none found"))
    }
}


func HaveNet6() {
    col := ColorBrush{enabled:!TerminalIsDumb()}
    ft := Formatter{colorBrush:col}

    // Detect local network info
    var ip6Ips = Ip6IpAddr()

    // Do local pings
//    var netPings = DoNetPings(ifconfig, route)

    // Detect ips on local area network
//    var lanIps = DetectLanIps(ifconfig, route)

    // Display local network info
    DisplayLocalNetwork6(ft, ip6Ips)
}
