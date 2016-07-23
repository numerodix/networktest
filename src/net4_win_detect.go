package main

//import "fmt"
//import "net"
//import "regexp"
//import "strings"


type WinNetDetect4 struct {
    ft Formatter
}


func WindowsNetworkDetector4(ft Formatter) WinNetDetect4 {
    return WinNetDetect4{
        ft: ft,
    }
}


func (wnd *WinNetDetect4) detectNetConn4() IP4NetworkInfo {
    var info = IP4NetworkInfo{}

    wnd.detectIpconfig4(&info)

    return info
}


func (wnd *WinNetDetect4) detectIpconfig4(info *IP4NetworkInfo) {
    var mgr = ProcMgr("ipconfig")
    var res = mgr.run()

    // The command failed :(
    if res.err != nil {
        wnd.ft.printError("Failed to detect ipv4 network", res.err)
        return
    }

    // Extract the output
    wnd.parseIpconfig4(res.stdout, info)

    // Parsing failed :(
    wnd.ft.printErrors("Failed to parse ipv4 network info", info.Errs)
}


func (wnd *WinNetDetect4) parseIpconfig4(stdout string, info *IP4NetworkInfo) {
}
