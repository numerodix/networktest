package main

import "fmt"
import "runtime"
import "strings"


type NetDetectUi struct {
    col ColorBrush
    ft Formatter

    ipver int
    osName string

    info4 *IP4NetworkInfo
}


func NetworkDetector(ipver int) NetDetectUi {
    var col = ColorBrush{
        enabled: !terminalIsDumb(),
    }
    var ft = Formatter{
        colorBrush: col,
    }

    var ui = NetDetectUi{
        col: col,
        ft: ft,
        ipver: ipver,
        osName: runtime.GOOS,
    }

    return ui
}


func (ui *NetDetectUi) run() {
    ui.displayPlatform()

    var info = ui.detectLocalNet()
    ui.info4 = &info

    ui.displayLocalNet()
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
        var winDet = WindowsNetworkDetector4(ui.ft)
        info = winDet.detectNetConn4()
    }

    info.normalize()

    return info
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
//        var pingExec = netPings[ifaceBlock.IPv4]
        var ifaceFmt = ui.ft.formatIfaceField(ip.Iface.Name)
        var ipFmt = ui.ft.formatIpField(ip.ipAsString())
        var maskFmt = ui.ft.formatSubnetField(ip.maskAsString())
//        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s %s   \n", ifaceFmt, ipFmt, maskFmt)
    }
    if len(ui.info4.Ips) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting gateways"))
    for _, gw := range ui.info4.getSortedGws() {
//        var pingExec = netPings[gw.Gateway]
        var ifaceFmt = ui.ft.formatIfaceField(gw.Iface.Name)
        var ipFmt = ui.ft.formatIpField(gw.ipAsString())
//        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s   \n", ifaceFmt, ipFmt)

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

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting dns servers"))
    for _, ip := range ui.info4.getSortedNsHosts() {
//        var pingExec = netPings[ifaceBlock.IPv4]
        var ipFmt = ui.ft.formatIpField(ip.ipAsString())
//        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s\n", ipFmt)
    }
    if len(ui.info4.NsHosts) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

}
