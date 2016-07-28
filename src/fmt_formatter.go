package main

import "fmt"
import "net"


type Formatter struct {
    colorBrush ColorBrush
}


func (ft *Formatter) formatPingTime(pingExec PingExecution) string {
    if pingExec.Unpingable {
        return ft.colorBrush.red("N/A")
    }

    // XXX: why is Err always nil?
    if pingExec.Err != nil || pingExec.Time <= 0.0 {
        return ft.colorBrush.red("failed")
    }

    var time = fmt.Sprintf("%.3f", pingExec.Time)
    time = time[:5]  // four significant digits + decimal point
    time = time + " ms"
    return ft.colorBrush.green(time)
}


func (ft *Formatter) formatHeader(msg string) string {
    var msgFmt = ft.colorBrush.yellow(fmt.Sprintf(" + %s...", msg))
    return msgFmt
}

func (ft *Formatter) formatError(msg string) string {
    var msgFmt = ft.colorBrush.red(fmt.Sprintf("%s", msg))
    return msgFmt
}

func (ft *Formatter) formatIfaceField(iface string) string {
    iface = fmt.Sprintf("<%s>", iface)
    var ifaceFmt = ft.colorBrush.magenta(fmt.Sprintf("%-10s", iface))
    return ifaceFmt
}

func (ft *Formatter) formatHostField(host string) string {
    var hostFmt = ft.colorBrush.cyan(fmt.Sprintf("%s", host))
    return hostFmt
}

func (ft *Formatter) formatInetHostField(host string) string {
    var hostFmt = ft.colorBrush.green(fmt.Sprintf("%15s", host))
    return hostFmt
}

func (ft *Formatter) formatLanIpField(ip net.IP) string {
    var ipFmt string

    if ipIs4(ip) {
        ipFmt = ft.colorBrush.bgreen(fmt.Sprintf("%15s", ip.String()))
    } else {
        ipFmt = ft.colorBrush.bgreen(fmt.Sprintf("%39s", ip.String()))
    }

    return ipFmt
}

func (ft *Formatter) formatIpField(ip net.IP) string {
    var ipFmt string

    if ipIs4(ip) {
        ipFmt = ft.colorBrush.green(fmt.Sprintf("%15s", ip.String()))
    } else {
        ipFmt = ft.colorBrush.green(fmt.Sprintf("%39s", ip.String()))
    }

    return ipFmt
}

func (ft *Formatter) formatIp6Field(ip net.IP) string {
    var ipFmt = ft.colorBrush.green(fmt.Sprintf("%39s", ip))
    return ipFmt
}

func (ft *Formatter) formatMask6Field(mask net.IPMask) string {
    var ones, _ = mask.Size()
    var subnetFmt = ft.colorBrush.cyan(fmt.Sprintf("/ %3d", ones))
    return subnetFmt
}

func (ft *Formatter) formatScope6Field(scope string) string {
    var ipFmt = ft.colorBrush.yellow(fmt.Sprintf("[scope: %s]", scope))
    return ipFmt
}

func (ft *Formatter) formatSubnetField(mask net.IPMask) string {
    var subnetFmt string
    var ones int

    if maskIs4(mask) {
        var maskStr = ipmaskAsString4(mask)
        subnetFmt = ft.colorBrush.cyan(fmt.Sprintf("/ %-15s", maskStr))
    } else {
        ones, _ = mask.Size()
        subnetFmt = ft.colorBrush.cyan(fmt.Sprintf("/ %3d", ones))
    }

    return subnetFmt
}


func (ft *Formatter) printError(msg string, err... error) {
    var prefix = fmt.Sprintf("Error: %s", msg)

    if err != nil {
        prefix = fmt.Sprintf("%s: %s", prefix, err)
    }

    fmt.Printf("%s\n", prefix)
}

func (ft *Formatter) printErrors(msg string, errs []error) {
    for _, err := range errs {
        ft.printError(msg, err)
    }
}
