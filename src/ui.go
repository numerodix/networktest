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

    var info = linuxDetectNetConn4()
    ui.displayLocalNet(&info)
}


func (ui *NetDetectUi) displayPlatform() {
    var plat = strings.Title(ui.osName)
    fmt.Printf("Platform: %s\n", ui.col.cyan(plat))
}


func (ui *NetDetectUi) displayLocalNet(info *IP4NetworkInfo) {

    fmt.Printf("%s\n", ui.ft.formatHeader("Scanning for networks"))
    for _, net := range info.Nets {
        var ifaceFmt = ui.ft.formatIfaceField(net.Iface.Name)
        var netwFmt = ui.ft.formatIpField(net.ipAsString())
        var maskFmt = ui.ft.formatSubnetField(net.maskAsString())
        fmt.Printf("    %s  %s %s\n", ifaceFmt, netwFmt, maskFmt)
    }
    if len(info.Nets) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting ips"))
    for _, ip := range info.Ips {
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
    for _, gw := range info.Gws {
//        var pingExec = netPings[gw.Gateway]
        var ifaceFmt = ui.ft.formatIfaceField(gw.Iface.Name)
        var ipFmt = ui.ft.formatIpField(gw.ipAsString())
//        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s   \n", ifaceFmt, ipFmt)
    }
/*    for _, lanIp := range lanIps {
        var ipFmt = ui.ft.formatLanIpField(lanIp)
        fmt.Printf("     ip:        %s\n", ipFmt)
    } */
    if len(info.Gws) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

}
