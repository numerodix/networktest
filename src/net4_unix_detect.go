package main

import "io/ioutil"
import "net"
import "regexp"
import "strings"


type UnixNetDetect4 struct {
    ft Formatter
}

func UnixNetworkDetector4(ft Formatter) UnixNetDetect4 {
    return UnixNetDetect4{
        ft: ft,
    }
}


func (und *UnixNetDetect4) unixDetectNsHosts4(info *IP4NetworkInfo) {
    var filepath = "/etc/resolv.conf"

    // Read the file
    var bytes, err = ioutil.ReadFile(filepath)
    if err != nil {
        und.ft.printError("Failed to detect ns servers", err)
        return
    }

    var content = string(bytes)
    und.unixParseResolvConf4(content, info)
}


func (und *UnixNetDetect4) unixParseResolvConf4(content string,
                                                info *IP4NetworkInfo) {

    var nameservers = und.unixParseResolvConf(content)

    for _, nameserver := range nameservers {
        var ip = net.ParseIP(nameserver)

        info.NsHosts = append(info.NsHosts, NsServer{
            Ip: ip,
        })
    }
}

func (und *UnixNetDetect4) unixParseResolvConf(content string) []string {
    // Parse the nameservers
    var nameservers []string
    var lines = strings.Split(content, "\n")
    rx := regexp.MustCompile("nameserver ([^ ]*)")

    for _, line := range lines {
        if rx.MatchString(line) {
            var ns = rx.FindStringSubmatch(line)[1]
            nameservers = append(nameservers, ns)
        }
    }

    return nameservers
}
