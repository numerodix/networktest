package main

import (
//    "bytes"
//    "errors"
    "fmt"
//    "log"
//    "os/exec"
//    "regexp"
//    "strconv"
//    "strings"
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


func main() {
    // Detect local network info
    var route = Route()
    var ifconfig = Ifconfig()


    // Do local pings
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


    // Do remote pings
    inetDnsServers := make(map[string]string)
    inetDnsServers["b.root-servers.net."] = "192.228.79.201"

    netDnsServers := []string{
        "127.0.1.1",
    }

    inetHosts := []string{
        "facebook.com",
        "github.com",
        "gmail.com",
        "google.com",
        "twitter.com",
        "nu.nl",
        "yahoo.com",
        "youtube.com",

        "aftenposten.no",
        "www.bonjourchine.com",
        "golang.org",
        "juventuz.com",
    }

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




    col := ColorBrush{enabled:true}
    ft := Formatter{colorBrush:col}

    /* 
        LOCAL NETWORK
    */

    fmt.Printf(col.yellow(" + Scanning for networks...\n"))
    var networks = route.GetNetworks()
    for i := range networks {
        var network = networks[i]

        var ifaceFmt = ft.FormatIfaceField(network.Iface)
        var netwFmt = ft.FormatIpField(network.Network)
        var maskFmt = ft.FormatSubnetField(network.Netmask)
        fmt.Printf("    %s  %s %s\n", ifaceFmt, netwFmt, maskFmt)
    }

    fmt.Printf(col.yellow(" + Detecting ips...\n"))
    for i := range ifaceBlocks {
        var ifaceBlock = ifaceBlocks[i]

        var pingExec = netPings[ifaceBlock.IPv4]
        var ifaceFmt = ft.FormatIfaceField(ifaceBlock.Iface)
        var ipFmt = ft.FormatIpField(ifaceBlock.IPv4)
        var maskFmt = ft.FormatSubnetField(ifaceBlock.Mask)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s  %s %s   ping: %s\n", ifaceFmt, ipFmt, maskFmt, pingFmt)
    }

    fmt.Printf(col.yellow(" + Detecting gateways...\n"))
    for i := range gws {
        var gw = gws[i]

        var pingExec = netPings[gw.Gateway]
        var ifaceFmt = ft.FormatIfaceField(gw.Iface)
        var ipFmt = ft.FormatIpField(gw.Gateway)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s  %s   ping: %s\n", ifaceFmt, ipFmt, pingFmt)
    }

    /* 
        INTERNET
    */

    fmt.Printf(col.yellow(" + Testing internet connection...\n"))
    for name, ip := range inetDnsServers {
        var pingExec = inetPings[ip]
        var nameFmt = ft.FormatHostField(name)
        var ipFmt = ft.FormatIpField(ip)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s  %s   ping: %s\n", nameFmt, ipFmt, pingFmt)
    }

    fmt.Printf(col.yellow(" + Detecting dns servers...\n"))
    for i := range netDnsServers {
        var host = netDnsServers[i]

        var pingExec = inetPings[host]
        var ipFmt = ft.FormatIpField(host)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s   ping: %s\n", ipFmt, pingFmt)
    }

    fmt.Printf(col.yellow(" + Testing internet dns...\n"))
    for i := range inetHosts {
        var host = inetHosts[i]

        var pingExec = inetPings[host]
        var ipFmt = ft.FormatIpField(host)
        var pingFmt = ft.FormatPingTime(pingExec)
        fmt.Printf("    %s   ping: %s\n", ipFmt, pingFmt)
    }
}
