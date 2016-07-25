package main

import "fmt"
import "runtime"
import "strings"


type NetDetectUi struct {
    ctx AppContext

    col ColorBrush
    ft Formatter

    ipver int  // 4 | 6
    osName string  // linux | freesd | ...

    // Well known nameservers on the internet (known by ip)
    inetNsHosts map[string]string
    // Well known hosts on the internet (known by hostname)
    inetHosts []string

    info4 *IP4NetworkInfo

    localPings Pings  // result of pinging local net
    inetPings Pings  // result of pinging inet
}


func NetworkDetector(ipver int) NetDetectUi {
    var col = ColorBrush{
        enabled: !terminalIsDumb(),
    }
    var ft = Formatter{
        colorBrush: col,
    }

    var ctx = AppContext{
        col: col,
        ft: ft,
        ipver: ipver,
        osName: runtime.GOOS,
    }


    var inetNsHosts = make(map[string]string)
    inetNsHosts["b.root-servers.net."] = "192.228.79.201"

    var inetHosts = []string{
        "facebook.com",
        "gmail.com",
        "google.com",
        "twitter.com",
        "yahoo.com",
    }

    var ui = NetDetectUi{
        ctx: ctx,
        col: col,
        ft: ft,
        ipver: ipver,
        osName: runtime.GOOS,

        inetNsHosts: inetNsHosts,
        inetHosts: inetHosts,
    }

    return ui
}


func (ui *NetDetectUi) run() {
    ui.displayPlatform()

    // Detect local network
    var info = ui.detectLocalNet()
    ui.info4 = &info

    // Ping local network
    ui.pingLocalNet()

    // Display local network
    ui.displayLocalNet()

    // If we don't have a local network connection we stop here
    if !ui.info4.haveLocalNet() {
        return
    }

    // Ping inet
    ui.pingInet()

    // Display inet connectivity
    ui.displayInetConnectivity()
}


func (ui *NetDetectUi) detectLocalNet() IP4NetworkInfo {
    var info IP4NetworkInfo

    switch ui.osName {
    // Linux userland
    case "linux":
        var linuxDet = LinuxNetworkDetector4(ui.ft)
        info = linuxDet.detectNetConn4()

    // BSD userland
    case "darwin":
        fallthrough
    case "dragonfly":
        fallthrough
    case "freebsd":
        fallthrough
    case "netbsd":
        fallthrough
    case "openbsd":
        var bsdDet = BsdNetworkDetector4(ui.ft)
        info = bsdDet.detectNetConn4()

    // Windows userland
    case "windows":
        var winDet = NewWinNetDetect4(ui.ctx)
        info = winDet.detectNetConn4()
    }

    info.normalize()

    return info
}


func (ui *NetDetectUi) pingLocalNet() {
    var hosts = []string{}

    // Ping local ips to see if reachable
    for _, ip := range ui.info4.Ips {
        hosts = append(hosts, ip.Ip.String())
    }

    // Ping gateways to see if reachable
    for _, gw := range ui.info4.Gws {
        hosts = append(hosts, gw.Ip.String())
    }

    // Run the pings
    var pinger = NewLinuxPinger4(ui.ft)
    var pings = runPings(pinger, hosts, 1, 1000)

    ui.localPings = pings
}


func (ui *NetDetectUi) pingInet() {
    var hosts = []string{}

    // Ping nameservers on inet to see if we can route packets there
    for _, ip := range ui.inetNsHosts {
        hosts = append(hosts, ip)
    }

    // Ping nameservers to see if we can resolve dns
    for _, nshost := range ui.info4.NsHosts {
        hosts = append(hosts, nshost.Ip.String())
    }

    // Ping hosts on inet to see if we can resolve dns and ping them
    for _, host := range ui.inetHosts {
        hosts = append(hosts, host)
    }

    // Run the pings
    var pinger = NewLinuxPinger4(ui.ft)
    var pings = runPings(pinger, hosts, 1, 2000)

    ui.inetPings = pings
}


func (ui *NetDetectUi) displayPlatform() {
    var plat = strings.Title(ui.osName)
    fmt.Printf("Platform: %s\n", ui.col.cyan(plat))
}


func (ui *NetDetectUi) displayLocalNet() {

    fmt.Printf("%s\n", ui.ft.formatHeader("Scanning for networks"))
    for _, net := range ui.info4.getSortedNets() {
        var ifaceFmt = ui.ft.formatIfaceField(net.Iface.Name)
        var netwFmt = ui.ft.formatIpField(net.ipAsString())
        var maskFmt = ui.ft.formatSubnetField(net.maskAsString())
        fmt.Printf("    %s  %s %s\n", ifaceFmt, netwFmt, maskFmt)
    }
    if len(ui.info4.Nets) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting ips"))
    for _, ip := range ui.info4.getSortedIps() {
        var pingExec = ui.localPings[ip.ipAsString()]
        var ifaceFmt = ui.ft.formatIfaceField(ip.Iface.Name)
        var ipFmt = ui.ft.formatIpField(ip.ipAsString())
        var maskFmt = ui.ft.formatSubnetField(ip.maskAsString())
        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s %s   ping: %s\n", ifaceFmt, ipFmt, maskFmt, pingFmt)
    }
    if len(ui.info4.Ips) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting gateways"))
    for _, gw := range ui.info4.getSortedGws() {
        var pingExec = ui.localPings[gw.ipAsString()]
        var ifaceFmt = ui.ft.formatIfaceField(gw.Iface.Name)
        var ipFmt = ui.ft.formatIpField(gw.ipAsString())
        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s   ping: %s\n", ifaceFmt, ipFmt, pingFmt)

        var ips = ui.info4.getIpsForGw(&gw)
        for _, ip := range ips {
            var ipFmt = ui.ft.formatLanIpField(ip.ipAsString())
            fmt.Printf("     ip:        %s\n", ipFmt)
        }
    }
    if len(ui.info4.Gws) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

}

func (ui *NetDetectUi) displayInetConnectivity() {

    fmt.Printf("%s\n", ui.ft.formatHeader("Testing internet connection"))
    for host, ip := range ui.inetNsHosts {
        var pingExec = ui.inetPings[ip]
        var nameFmt = ui.ft.formatHostField(host)
        var ipFmt = ui.ft.formatIpField(ip)
        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s  ping: %s\n", nameFmt, ipFmt, pingFmt)
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting dns servers"))
    for _, ip := range ui.info4.getSortedNsHosts() {
        var pingExec = ui.inetPings[ip.ipAsString()]
        var ipFmt = ui.ft.formatIpField(ip.ipAsString())
        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s   ping: %s\n", ipFmt, pingFmt)
    }
    if len(ui.info4.NsHosts) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Testing internet dns"))
    for _, host := range ui.inetHosts {
        var pingExec = ui.inetPings[host]
        var ipFmt = ui.ft.formatIpField(host)
        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s   ping: %s\n", ipFmt, pingFmt)
    }

}
