package main

import (
    "fmt"
)


func PingJob(ch chan PingExecution, host string, cnt int, timeout int) {
    ch <- Ping(host, cnt, timeout)
}


func SpawnAndCollect(pingHosts map[string]PingExecution) {
    ch := make(chan PingExecution)

    // Launch
    for ip := range pingHosts {
        go PingJob(ch, ip, 1, 2)
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


func DoNetPings(ifconfig IfconfigExecution, route RouteExecution) map[string]PingExecution {
    var netPings = make(map[string]PingExecution)

    var ifaceBlocks = ifconfig.IfaceBlocks
    for i := range ifaceBlocks {
        var ifaceBlock = ifaceBlocks[i]
        var ip = ifaceBlock.IPv4

        netPings[ip] = PingExecution{}
    }

    var gws = route.GetGateways()
    for i := range gws {
        var gw = gws[i]
        var ip = gw.Gateway

        netPings[ip] = PingExecution{}
    }

    SpawnAndCollect(netPings)

    return netPings
}


func DoInetPings(inetDnsServers map[string]string, netDnsServers []string,
                 inetHosts []string) map[string]PingExecution {
    var inetPings = make(map[string]PingExecution)

    for _, ip := range inetDnsServers {
        inetPings[ip] = PingExecution{}
    }

    for i := range netDnsServers {
        var host = netDnsServers[i]
        inetPings[host] = PingExecution{}
    }

    for i := range inetHosts {
        var host = inetHosts[i]
        inetPings[host] = PingExecution{}
    }

    SpawnAndCollect(inetPings)

    return inetPings
}


func DetectLanIps(ifconfig IfconfigExecution, route RouteExecution) []string {
    var gws = route.GetGateways()
    var ifaceBlocks = ifconfig.IfaceBlocks

    var lanIps []string

    for i := range gws {
        var gw = gws[i]

        for j := range ifaceBlocks {
            var ifaceBlock = ifaceBlocks[j]

            var gwIp = gw.Gateway
            var hostIp = ifaceBlock.IPv4
            var maskIp = ifaceBlock.Mask

            if NetworkContains(hostIp, maskIp, gwIp) {
                lanIps = append(lanIps, hostIp)
            }
        }
    }

    return lanIps
}


func DisplayLocalNetwork(ft Formatter,
                         ifconfig IfconfigExecution,
                         route RouteExecution,
                         lanIps []string,
                         netPings map[string]PingExecution) {

    var gws = route.GetGateways()
    var ifaceBlocks = ifconfig.IfaceBlocks

    fmt.Printf("%s\n", ft.FormatHeader("Scanning for networks"))
    var networks = route.GetNetworks()
    for i := range networks {
        var network = networks[i]

        var ifaceFmt = ft.FormatIfaceField(network.Iface)
        var netwFmt = ft.FormatIpField(network.Network)
        var maskFmt = ft.FormatSubnetField(network.Netmask)
        fmt.Printf("    %s  %s %s\n", ifaceFmt, netwFmt, maskFmt)
    }
    if len(networks) == 0 {
        fmt.Printf("    %s\n", ft.FormatError("none found"))
    }

    fmt.Printf("%s\n", ft.FormatHeader("Detecting ips"))
    for i := range ifaceBlocks {
        var ifaceBlock = ifaceBlocks[i]

        var pingExec = netPings[ifaceBlock.IPv4]
        var ifaceFmt = ft.FormatIfaceField(ifaceBlock.Iface)
        var ipFmt = ft.FormatIpField(ifaceBlock.IPv4)
        var maskFmt = ft.FormatSubnetField(ifaceBlock.Mask)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s  %s %s   ping: %s\n", ifaceFmt, ipFmt, maskFmt, pingFmt)
    }
    if len(ifaceBlocks) == 0 {
        fmt.Printf("    %s\n", ft.FormatError("none found"))
    }

    fmt.Printf("%s\n", ft.FormatHeader("Detecting gateways"))
    for i := range gws {
        var gw = gws[i]

        var pingExec = netPings[gw.Gateway]
        var ifaceFmt = ft.FormatIfaceField(gw.Iface)
        var ipFmt = ft.FormatIpField(gw.Gateway)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s  %s   ping: %s\n", ifaceFmt, ipFmt, pingFmt)
    }
    for i := range lanIps {
        var lanIp = lanIps[i]

        var ipFmt = ft.FormatLanIpField(lanIp)
        fmt.Printf("     ip:        %s\n", ipFmt)
    }
    if len(gws) == 0 {
        fmt.Printf("    %s\n", ft.FormatError("none found"))
    }
}


func DisplayInetConnectivity(ft Formatter,
                             inetDnsServers map[string]string, netDnsServers []string,
                             inetHosts []string,
                             inetPings map[string]PingExecution) {

    fmt.Printf("%s\n", ft.FormatHeader("Testing internet connection"))
    for name, ip := range inetDnsServers {
        var pingExec = inetPings[ip]
        var nameFmt = ft.FormatHostField(name)
        var ipFmt = ft.FormatIpField(ip)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s  %s  ping: %s\n", nameFmt, ipFmt, pingFmt)
    }

    fmt.Printf("%s\n", ft.FormatHeader("Detecting dns servers"))
    for i := range netDnsServers {
        var host = netDnsServers[i]

        var pingExec = inetPings[host]
        var ipFmt = ft.FormatIpField(host)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s   ping: %s\n", ipFmt, pingFmt)
    }

    fmt.Printf("%s\n", ft.FormatHeader("Testing internet dns"))
    for i := range inetHosts {
        var host = inetHosts[i]

        var pingExec = inetPings[host]
        var ipFmt = ft.FormatIpField(host)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s   ping: %s\n", ipFmt, pingFmt)
    }
}


func main() {
    col := ColorBrush{enabled:!TerminalIsDumb()}
    ft := Formatter{colorBrush:col}

    // Detect local network info
    var route = Route()
    var ifconfig = Ifconfig()

    // Do local pings
    var netPings = DoNetPings(ifconfig, route)

    // Detect ips on local area network
    var lanIps = DetectLanIps(ifconfig, route)

    // Display local network info
    DisplayLocalNetwork(ft, ifconfig, route, lanIps, netPings)

    // If we didn't find any lan ips, don't test inet connectivity
    if len(lanIps) == 0 {
        return
    }

    // Do remote pings
    inetDnsServers := make(map[string]string)
    inetDnsServers["b.root-servers.net."] = "192.228.79.201"

    var netDnsServers = DetectNameservers()

    inetHosts := []string{
        "facebook.com",
        "gmail.com",
        "google.com",
        "twitter.com",
        "yahoo.com",
    }

    // Do inet pings
    var inetPings = DoInetPings(inetDnsServers, netDnsServers, inetHosts)

    // Display inet connectivity info
    DisplayInetConnectivity(ft, inetDnsServers, netDnsServers, inetHosts, inetPings)
}
