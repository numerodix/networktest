package main

import "fmt"
import "runtime"


type NetDetectUi struct {
    col ColorBrush
    ft Formatter
    osName string
}


func NetworkDetector() NetDetectUi {
    var col = ColorBrush{
        enabled: !TerminalIsDumb(),
    }
    var ft = Formatter{
        colorBrush: col,
    }

    var ui = NetDetectUi{
        col: col,
        ft: ft,
        osName: runtime.GOOS,
    }

    return ui
}


func run(ui *NetDetectUi) {
    var info = linuxDetectNetConn4()

    displayLocalNet(ui, &info)
}


func displayLocalNet(ui *NetDetectUi, info *IP4NetworkInfo) {

    fmt.Printf("%s\n", ui.ft.formatHeader("Scanning for networks"))
    for i := range info.Nets {
        var net = info.Nets[i]

        var ifaceFmt = ui.ft.formatIfaceField(net.Iface.Name)
        var netwFmt = ui.ft.formatIpField(net.Ip.IP.String())
        var maskFmt = ui.ft.formatSubnetField(net.Ip.Mask.String())
        fmt.Printf("    %s  %s %s\n", ifaceFmt, netwFmt, maskFmt)
    }
    if len(info.Nets) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting ips"))
    for i := range info.Ips {
        var ip = info.Ips[i]

//        var pingExec = netPings[ifaceBlock.IPv4]
        var ifaceFmt = ui.ft.formatIfaceField(ip.Iface.Name)
        var ipFmt = ui.ft.formatIpField(ip.Ip.String())
        var maskFmt = ui.ft.formatSubnetField(ip.Mask.String())
//        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s %s   \n", ifaceFmt, ipFmt, maskFmt)
    }
    if len(info.Ips) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

    fmt.Printf("%s\n", ui.ft.formatHeader("Detecting gateways"))
    for i := range info.Gws {
        var gw = info.Gws[i]

//        var pingExec = netPings[gw.Gateway]
        var ifaceFmt = ui.ft.formatIfaceField(gw.Iface.Name)
        var ipFmt = ui.ft.formatIpField(gw.Ip.String())
//        var pingFmt = ui.ft.formatPingTime(pingExec)
        fmt.Printf("    %s  %s   \n", ifaceFmt, ipFmt)
    }
/*    for i := range lanIps {
        var lanIp = lanIps[i]

        var ipFmt = ui.ft.formatLanIpField(lanIp)
        fmt.Printf("     ip:        %s\n", ipFmt)
    } */
    if len(info.Gws) == 0 {
        fmt.Printf("    %s\n", ui.ft.formatError("none found"))
    }

}
