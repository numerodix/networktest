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


type CommandResult struct {
    Id string
    FValue float64
    SValue string
    Error error
}


func PingProc(ch chan CommandResult, host string, cnt int, timeout int) {
    pingExec := Ping(host, cnt, timeout)
    if pingExec.Error != nil {
        ch <- CommandResult{Id: pingExec.Host, Error: pingExec.Error}
        return
    }

    ch <- CommandResult{Id: pingExec.Host, FValue: pingExec.Time}
}


func main() {
//    hosts := []string{"yahoo.com", "google.com"}
    hosts := []string{
        "192.168.1.1",
/*        "192.228.79.201",
        "127.0.0.1",
        "127.0.1.1",
        "localhost",
        "yahoo.com",
        "google.com",
        "juventuz.com",
        "twitter.com",
        "facebook.com",
        "gmail.com",
        "golang.org",
        "www.nu.nl",
        "www.aftenposten.no",
        "www.bonjourchine.com",
        "github.com",
        "youtube.com",
*/
    }
//    hosts := []string{"localhost"}
    ch := make(chan CommandResult)

    // Launch
    for i := range hosts {
        go PingProc(ch, hosts[i], 1, 2)
    }

    // Collect
    sum := 0.0
    for i := range hosts {
        cmdres := <-ch

        if cmdres.Error != nil {
            fmt.Printf("Err: %s: %s\n", cmdres.Id, cmdres.Error)
            continue
        }

        sum += cmdres.FValue
        fmt.Printf("%-2d  %-34s: %.1f ms\n", i, cmdres.Id, cmdres.FValue)
    }

    fmt.Printf("Total time: %.1f ms\n", sum)

    fmt.Printf("ping exec: %s\n", Ping("localhost", 1, 1))


    var route = Route()
    fmt.Printf("networks: %s\n", route.GetNetworks())
    fmt.Printf("gateways: %s\n", route.GetGateways())

    var ifconfig = Ifconfig()
    fmt.Printf("ifconfig: %s\n", ifconfig)



    // Do local pings
    var localPings = make(map[string]PingExecution)

    var ifaceBlocks = ifconfig.IfaceBlocks
    for i := range ifaceBlocks {
        var ifaceBlock = ifaceBlocks[i]
        var ip = ifaceBlock.IPv4

        localPings[ip] = Ping(ip, 1, 2)
    }

    var gws = route.GetGateways()
    for i := range gws {
        var gw = gws[i]
        var ip = gw.Gateway

        localPings[ip] = Ping(ip, 1, 2)
    }





    col := ColorBrush{enabled:true}

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
//    var ifaceBlocks = ifconfig.IfaceBlocks
    for i := range ifaceBlocks {
        var ifaceBlock = ifaceBlocks[i]

        var iface = fmt.Sprintf("<%s>", ifaceBlock.Iface)
        var ip = ifaceBlock.IPv4
        var mask = ifaceBlock.Mask
        var ping = localPings[ip].Time
        var ifaceS = col.magenta(fmt.Sprintf("%-10s", iface))
        var ipS = col.green(fmt.Sprintf("%-15s", ip))
        var maskS = col.cyan(fmt.Sprintf("/ %-15s", mask))
        var pingS = col.green(fmt.Sprintf("%.3f ms", ping))
        fmt.Printf("    %s  %s %s   ping: %s\n", ifaceS, ipS, maskS, pingS)
    }

    fmt.Printf(col.yellow(" + Detecting gateways...\n"))
//    var gws = route.GetGateways()
    for i := range gws {
        var gw = gws[i]

        var iface = fmt.Sprintf("<%s>", gw.Iface)
        var ip = gw.Gateway
        var ping = localPings[ip].Time
        var ifaceS = col.magenta(fmt.Sprintf("%-10s", iface))
        var ipS = col.green(fmt.Sprintf("%-15s", ip))
        var pingS = col.green(fmt.Sprintf("%.3f ms", ping))
        fmt.Printf("    %s  %s   ping: %s\n", ifaceS, ipS, pingS)
    }

}
