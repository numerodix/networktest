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

    /* 
        LOCAL NETWORK
    */

    fmt.Printf(col.yellow(" + Scanning for networks...\n"))
    var networks = route.GetNetworks()
    for i := range networks {
        var network = networks[i]

        var iface = fmt.Sprintf("<%s>", network.Iface)
        var netw = network.Network
        var mask = network.Netmask
        var ifaceS = col.magenta(fmt.Sprintf("%-10s", iface))
        var netwS = col.green(fmt.Sprintf("%-15s", netw))
        var maskS = col.cyan(fmt.Sprintf("/ %s", mask))
        fmt.Printf("    %s  %s %s\n", ifaceS, netwS, maskS)
    }

    fmt.Printf(col.yellow(" + Detecting ips...\n"))
    for i := range ifaceBlocks {
        var ifaceBlock = ifaceBlocks[i]

        var iface = fmt.Sprintf("<%s>", ifaceBlock.Iface)
        var ip = ifaceBlock.IPv4
        var mask = ifaceBlock.Mask
        var ping = netPings[ip].Time
        var ifaceS = col.magenta(fmt.Sprintf("%-10s", iface))
        var ipS = col.green(fmt.Sprintf("%-15s", ip))
        var maskS = col.cyan(fmt.Sprintf("/ %-15s", mask))
        var pingS = col.green(fmt.Sprintf("%.3f ms", ping))
        fmt.Printf("    %s  %s %s   ping: %s\n", ifaceS, ipS, maskS, pingS)
    }

    fmt.Printf(col.yellow(" + Detecting gateways...\n"))
    for i := range gws {
        var gw = gws[i]

        var iface = fmt.Sprintf("<%s>", gw.Iface)
        var ip = gw.Gateway
        var ping = netPings[ip].Time
        var ifaceS = col.magenta(fmt.Sprintf("%-10s", iface))
        var ipS = col.green(fmt.Sprintf("%-15s", ip))
        var pingS = col.green(fmt.Sprintf("%.3f ms", ping))
        fmt.Printf("    %s  %s   ping: %s\n", ifaceS, ipS, pingS)
    }

    /* 
        INTERNET
    */

    fmt.Printf(col.yellow(" + Testing internet connection...\n"))
    for name, ip := range inetDnsServers {
        var ping = inetPings[ip].Time
        var nameS = col.cyan(fmt.Sprintf("%s", name))
        var ipS = col.green(fmt.Sprintf("%-15s", ip))
        var pingS = col.green(fmt.Sprintf("%.1f ms", ping))
        fmt.Printf("    %s  %s   ping: %s\n", nameS, ipS, pingS)
    }

    fmt.Printf(col.yellow(" + Detecting dns servers...\n"))
    for i := range netDnsServers {
        var host = netDnsServers[i]

        var ip = host
        var ping = inetPings[host].Time
        var ipS = col.green(fmt.Sprintf("%-15s", ip))
        var pingS = col.green(fmt.Sprintf("%.1f ms", ping))
        fmt.Printf("    %s   ping: %s\n", ipS, pingS)
    }

    fmt.Printf(col.yellow(" + Testing internet dns...\n"))
    for i := range inetHosts {
        var host = inetHosts[i]

        var ip = host
        var ping = inetPings[host].Time
        var ipS = col.green(fmt.Sprintf("%-15s", ip))
        var pingS = col.green(fmt.Sprintf("%.1f ms", ping))
        fmt.Printf("    %s   ping: %s\n", ipS, pingS)
    }
}
