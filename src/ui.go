package main

import "fmt"
import "runtime"
import "strings"


type NetDetectUi struct {
    col ColorBrush
    ft Formatter
    ipver int
    osName string
}


func NetworkDetector(ipver int) NetDetectUi {
    var col = ColorBrush{
        enabled: !TerminalIsDumb(),
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

    var info IP4NetworkInfo

    switch ui.osName {
    // Linux userland
    case "linux":
        var linuxDet = LinuxNetworkDetector4(ui.ft)
        info = linuxDet.detectNetConn4()
        break

    // BSD userland
    case "darwin":
    case "dragonfly":
    case "freebsd":
    case "netbsd":
    case "openbsd":
        var bsdDet = BsdNetworkDetector4(ui.ft)
        info = bsdDet.detectNetConn4()
        break
    }

    ui.displayLocalNet(&info)
}


func (ui *NetDetectUi) displayPlatform() {
    var plat = strings.Title(ui.osName)
    fmt.Printf("Platform: %s\n", ui.col.cyan(plat))
}


func (ui *NetDetectUi) displayLocalNet(info *IP4NetworkInfo) {

    fmt.Printf("%s\n", ui.ft.formatHeader("Scanning for networks"))
    for _, net := range info.getSortedNets() {
        var ifaceFmt = ui.ft.formatIfaceField(net.Iface.Name)
        var netwFmt = ui.ft.formatIpField(net.ipAsString())
        var maskFmt = ui.ft.formatSubnetField(net.maskAsString())
        fmt.Printf("    %s  %s %s\n", ifaceFmt, netwFmt, maskFmt)
    }
    if len(info.Nets) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting ips"))
    for _, ip := range info.getSortedIps() {
//        var pingExec = netPings[ifaceBlock.IPv4]
        var ifaceFmt = ui.ft.formatIfaceField(ip.Iface.Name)
        var ipFmt = ui.ft.formatIpField(ip.ipAsString())
        var maskFmt = ui.ft.formatSubnetField(ip.maskAsString())
//        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s %s   \n", ifaceFmt, ipFmt, maskFmt)
    }
    if len(info.Ips) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting gateways"))
    for _, gw := range info.getSortedGws() {
//        var pingExec = netPings[gw.Gateway]
        var ifaceFmt = ui.ft.formatIfaceField(gw.Iface.Name)
        var ipFmt = ui.ft.formatIpField(gw.ipAsString())
//        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s   \n", ifaceFmt, ipFmt)

        var ips = info.getIpsForGw(&gw)
        for _, ip := range ips {
            var ipFmt = ui.ft.formatLanIpField(ip.ipAsString())
            fmt.Printf("     ip:        %s\n", ipFmt)
        }
    }
    if len(info.Gws) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

}
