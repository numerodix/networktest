package main

import "fmt"
import "net"
import "runtime"
import "strings"


type NetDetectUi struct {
    ctx AppContext

    // Well known nameservers on the internet (known by ip)
    inet4NsHosts map[string]net.IP
    inet6NsHosts map[string]net.IP
    // Well known hosts on the internet (known by hostname)
    inetHosts []string

    info *IPNetworkInfo

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


    var inet4NsHosts = make(map[string]net.IP)
    inet4NsHosts["b.root-servers.net."] = net.ParseIP("192.228.79.201")

    var inet6NsHosts = make(map[string]net.IP)
    inet6NsHosts["b.root-servers.net."] = net.ParseIP("2001:500:84::b")

    var inetHosts = []string{
        "facebook.com",
        "gmail.com",
        "google.com",
        "twitter.com",
        "yahoo.com",
    }


    var ui = NetDetectUi{
        ctx: ctx,

        inet4NsHosts: inet4NsHosts,
        inet6NsHosts: inet6NsHosts,
        inetHosts: inetHosts,
    }

    return ui
}


func (ui *NetDetectUi) getInetNsHosts() map[string]net.IP {
    if ui.ctx.ipver == 4 {
        return ui.inet4NsHosts
    } else {
        return ui.inet6NsHosts
    }
}


func (ui *NetDetectUi) run() {
    ui.displayPlatform()

    // Detect local network
    var info = ui.detectLocalNet()
    ui.info = &info

    // Ping local network
    ui.pingLocalNet()

    // Display local network
    ui.displayLocalNet()

    // If we don't have a local network connection we stop here
    if !ui.info.haveLocalNet() {
        return
    }

    // Ping inet
    ui.pingInet()

    // Display inet connectivity
    ui.displayInetConnectivity()
}


func (ui *NetDetectUi) detectLocalNet() IPNetworkInfo {
    var info IPNetworkInfo

    switch ui.ctx.ipver {
    case 4:
        var detector = getDetector4(ui.ctx)
        if detector != nil {  // in case platform unsupported
            info = detector.detectNetConn4()
        }

    case 6:
        var detector = getDetector6(ui.ctx)
        if detector != nil {  // in case platform unsupported
            info = detector.detectNetConn6()
        }
    }

    info.normalize()
    return info
}


func (ui *NetDetectUi) pingLocalNet() {
    var hosts = []string{}

    // Ping local ips to see if reachable
    for _, ip := range ui.info.Ips {
        hosts = append(hosts, ip.Ip.String())
    }

    // Ping gateways to see if reachable
    for _, gw := range ui.info.Gws {
        hosts = append(hosts, gw.Ip.String())
    }

    // Run the pings
    var pinger = getPinger(ui.ctx)
    var pings = runPings(pinger, hosts, 1, 1000)

    ui.localPings = pings
}


func (ui *NetDetectUi) pingInet() {
    var hosts = []string{}

    // Ping nameservers on inet to see if we can route packets there
    for _, ip := range ui.getInetNsHosts() {
        hosts = append(hosts, ip.String())
    }

    // Ping nameservers to see if we can resolve dns
    for _, nshost := range ui.info.NsHosts {
        hosts = append(hosts, nshost.Ip.String())
    }

    // Ping hosts on inet to see if we can resolve dns and ping them
    for _, host := range ui.inetHosts {
        hosts = append(hosts, host)
    }

    // Run the pings
    var pinger = getPinger(ui.ctx)
    var pings = runPings(pinger, hosts, 1, 2000)

    ui.inetPings = pings
}


func (ui *NetDetectUi) displayPlatform() {
    var plat = strings.Title(ui.ctx.osName)
    fmt.Printf("Platform: %s\n", ui.ctx.col.cyan(plat))
}


func (ui *NetDetectUi) displayLocalNet() {

    fmt.Printf("%s\n", ui.ctx.ft.formatHeader("Scanning for networks"))
    for _, net := range ui.info.getSortedNets() {
        var ifaceFmt = ui.ctx.ft.formatIfaceField(net.Iface.Name)
        var netwFmt = ui.ctx.ft.formatIpField(net.Ip.IP)
        var maskFmt = ui.ctx.ft.formatSubnetField(net.Ip.Mask)

        var scopeFmt = ""
        if ipIs6(net.Ip.IP) {
            var scope = ip6AsScope(net.Ip.IP)
            scopeFmt = ui.ctx.ft.formatScope6Field(scope)
        }

        fmt.Printf("    %s  %s %s  %s\n", ifaceFmt, netwFmt, maskFmt, scopeFmt)
    }
    if len(ui.info.Nets) == 0 {
        fmt.Printf("    %s\n", ui.ctx.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ctx.ft.formatHeader("Detecting ips"))
    for _, ip := range ui.info.getSortedIps() {
        var pingExec = ui.localPings[ip.ipAsString()]
        var ifaceFmt = ui.ctx.ft.formatIfaceField(ip.Iface.Name)
        var ipFmt = ui.ctx.ft.formatIpField(ip.Ip)
        var maskFmt = ui.ctx.ft.formatSubnetField(ip.maskAsIPMask())
        var pingFmt = ui.ctx.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s %s  ping: %s\n", ifaceFmt, ipFmt, maskFmt, pingFmt)
    }
    if len(ui.info.Ips) == 0 {
        fmt.Printf("    %s\n", ui.ctx.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ctx.ft.formatHeader("Detecting gateways"))
    for _, gw := range ui.info.getSortedGws() {
        var pingExec = ui.localPings[gw.ipAsString()]
        var ifaceFmt = ui.ctx.ft.formatIfaceField(gw.Iface.Name)
        var ipFmt = ui.ctx.ft.formatIpField(gw.Ip)
        var pingFmt = ui.ctx.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s  ping: %s\n", ifaceFmt, ipFmt, pingFmt)

        var ips = ui.info.getIpsForGw(&gw)
        for _, ip := range ips {
            var ipFmt = ui.ctx.ft.formatLanIpField(ip.Ip)
            fmt.Printf("     ip:        %s\n", ipFmt)
        }
    }
    if len(ui.info.Gws) == 0 {
        fmt.Printf("    %s\n", ui.ctx.ft.formatError("none found"))
    }

}

func (ui *NetDetectUi) displayInetConnectivity() {

    fmt.Printf("%s\n", ui.ctx.ft.formatHeader("Testing internet connection"))
    for host, ip := range ui.getInetNsHosts() {
        var pingExec = ui.inetPings[ip.String()]
        var nameFmt = ui.ctx.ft.formatHostField(host)
        var ipFmt = ui.ctx.ft.formatIpField(ip)
        var pingFmt = ui.ctx.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s  ping: %s\n", nameFmt, ipFmt, pingFmt)
    }

    fmt.Printf("%s\n", ui.ctx.ft.formatHeader("Detecting dns servers"))
    for _, ip := range ui.info.getSortedNsHosts() {
        var pingExec = ui.inetPings[ip.ipAsString()]
        var ipFmt = ui.ctx.ft.formatIpField(ip.Ip)
        var pingFmt = ui.ctx.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  ping: %s\n", ipFmt, pingFmt)
    }
    if len(ui.info.NsHosts) == 0 {
        fmt.Printf("    %s\n", ui.ctx.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ctx.ft.formatHeader("Testing internet dns"))
    for _, host := range ui.inetHosts {
        var pingExec = ui.inetPings[host]
        var ipFmt = ui.ctx.ft.formatInetHostField(host)
        var pingFmt = ui.ctx.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  ping: %s\n", ipFmt, pingFmt)
    }

}
