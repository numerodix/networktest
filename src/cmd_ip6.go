package main

import (
    "fmt"
    "net"
)


func PingJob6(ch chan PingExecution, host string, cnt int, timeout int) {
    var pinger = newPinger6()
    ch <- pinger.Ping(host, cnt, timeout)
}


func SpawnAndCollect6(pingHosts map[string]PingExecution) {
    ch := make(chan PingExecution)

    // Launch
    for ip := range pingHosts {
        go PingJob6(ch, ip, 1, 2)
    }

    // Collect
    for ip := range pingHosts {
        pingExec := <-ch

        pingHosts[pingExec.Host] = pingExec

        // To make it not complain about ip not being used lol
        if false {
            fmt.Printf("%s", ip)
        }
    }
}


func DoNetPings6(ip6Addrs Ip6AddrExecution,
                 ip6Routes Ip6RouteExecution) map[string]PingExecution {

    var netPings = make(map[string]PingExecution)

    var ip6AddrBlocks = ip6Addrs.Ip6AddrBlocks
    for i := range ip6AddrBlocks {
        var ip6AddrBlock = ip6AddrBlocks[i]
        var ip = ip6AddrBlock.IPv6.String()

        netPings[ip] = PingExecution{}
    }

    var ip6RouteBlocks = ip6Routes.Ip6RouteBlocks
    for i := range ip6RouteBlocks {
        var ip6RouteBlock = ip6RouteBlocks[i]
        var ip = ip6RouteBlock.IPv6.String()

        netPings[ip] = PingExecution{}
    }

    SpawnAndCollect6(netPings)

    return netPings
}


func DetectLanIps6(ip6Addrs Ip6AddrExecution,
                   ip6Routes Ip6RouteExecution) []net.IP {

    var ip6AddrBlocks = ip6Addrs.Ip6AddrBlocks
    var ip6RouteBlocks = ip6Routes.Ip6RouteBlocks

    var lanIps []net.IP

    for i := range ip6RouteBlocks {
        var ip6RouteBlock = ip6RouteBlocks[i]

        for j := range ip6AddrBlocks {
            var ip6AddrBlock = ip6AddrBlocks[j]

            if ip6AddrBlock.Network.Contains(ip6RouteBlock.IPv6) {
                lanIps = append(lanIps, ip6AddrBlock.IPv6)
            }
        }
    }

    return lanIps
}


func DisplayLocalNetwork6(ft Formatter,
                          ip6Addrs Ip6AddrExecution,
                          ip6Routes Ip6RouteExecution,
                          lanIps []net.IP,
                          netPings map[string]PingExecution) {

    var ip6AddrBlocks = ip6Addrs.Ip6AddrBlocks
    var ip6RouteBlocks = ip6Routes.Ip6RouteBlocks

    fmt.Printf("%s\n", ft.FormatHeader("Scanning for networks"))
    for i := range ip6AddrBlocks {
        var ip6AddrBlock = ip6AddrBlocks[i]

        var ifaceFmt = ft.FormatIfaceField(ip6AddrBlock.Iface)
        var netwFmt = ft.FormatIp6Field(ip6AddrBlock.Network.IP)
        var maskFmt = ft.FormatMask6Field(ip6AddrBlock.Network.Mask)
        var scopeFmt = ft.FormatScope6Field(ip6AddrBlock.Scope)
        fmt.Printf("    %s  %s %s  %s\n", ifaceFmt, netwFmt, maskFmt, scopeFmt)
    }
    if len(ip6AddrBlocks) == 0 {
        fmt.Printf("    %s\n", ft.FormatError("none found"))
    }

    fmt.Printf("%s\n", ft.FormatHeader("Detecting ips"))
    for i := range ip6AddrBlocks {
        var ip6AddrBlock = ip6AddrBlocks[i]

        var pingExec = netPings[ip6AddrBlock.IPv6.String()]
        var ifaceFmt = ft.FormatIfaceField(ip6AddrBlock.Iface)
        var ipFmt = ft.FormatIp6Field(ip6AddrBlock.IPv6)
        var maskFmt = ft.FormatMask6Field(ip6AddrBlock.Network.Mask)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s  %s %s  ping: %s\n", ifaceFmt, ipFmt, maskFmt, pingFmt)
    }
    if len(ip6AddrBlocks) == 0 {
        fmt.Printf("    %s\n", ft.FormatError("none found"))
    }

    fmt.Printf("%s\n", ft.FormatHeader("Detecting gateways"))
    for i := range ip6RouteBlocks {
        var ip6RouteBlock = ip6RouteBlocks[i]

        var pingExec = netPings[ip6RouteBlock.IPv6.String()]
        var ifaceFmt = ft.FormatIfaceField(ip6RouteBlock.Iface)
        var ipFmt = ft.FormatIp6Field(ip6RouteBlock.IPv6)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s  %s   ping: %s\n", ifaceFmt, ipFmt, pingFmt)
    }
    for i := range lanIps {
        var lanIp = lanIps[i]

        var ipFmt = ft.FormatLanIp6Field(lanIp)
        fmt.Printf("     ip:        %s\n", ipFmt)
    }
    if len(ip6RouteBlocks) == 0 {
        fmt.Printf("    %s\n", ft.FormatError("none found"))
    }
}


func HaveNet6() {
    col := ColorBrush{enabled:!TerminalIsDumb()}
    ft := Formatter{colorBrush:col}

    // Detect local network info
    var ip6Addrs = Ip6Addr()
    var ip6Routes = Ip6Route()

    // Do local pings
    var netPings = DoNetPings6(ip6Addrs, ip6Routes)

    // Detect ips on local area network
    var lanIps = DetectLanIps6(ip6Addrs, ip6Routes)

    // Display local network info
    DisplayLocalNetwork6(ft, ip6Addrs, ip6Routes, lanIps, netPings)
}
