package main

import (
    "fmt"
)


func DisplayLocalNetwork6(ft Formatter,
                          ip6Ips Ip6AddrExecution) {

    var ip6AddrBlocks = ip6Ips.Ip6AddrBlocks

    fmt.Printf("%s\n", ft.FormatHeader("Scanning for networks"))
    for i := range ip6AddrBlocks {
        var ip6AddrBlock = ip6AddrBlocks[i]

        var ifaceFmt = ft.FormatIfaceField(ip6AddrBlock.Iface)
        var netwFmt = ft.FormatIp6Field(ip6AddrBlock.Network.IP)
        var maskFmt = ft.FormatMask6Field(ip6AddrBlock.Network.Mask)
        var scopeFmt = ft.FormatScope6Field(ip6AddrBlock.Scope)
        fmt.Printf("    %s  %s %s   %s\n", ifaceFmt, netwFmt, maskFmt, scopeFmt)
    }
    if len(ip6AddrBlocks) == 0 {
        fmt.Printf("    %s\n", ft.FormatError("none found"))
    }

    fmt.Printf("%s\n", ft.FormatHeader("Detecting ips"))
    for i := range ip6AddrBlocks {
        var ip6AddrBlock = ip6AddrBlocks[i]

        var ifaceFmt = ft.FormatIfaceField(ip6AddrBlock.Iface)
        var ipFmt = ft.FormatIp6Field(ip6AddrBlock.IPv6)
        var maskFmt = ft.FormatMask6Field(ip6AddrBlock.Network.Mask)
        fmt.Printf("    %s  %s %s\n", ifaceFmt, ipFmt, maskFmt)
    }
    if len(ip6AddrBlocks) == 0 {
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
