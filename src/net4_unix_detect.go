package main

import "io/ioutil"
import "net"
import "regexp"
import "strings"


type UnixNetDetect4 struct {
    ctx AppContext
}

func NewUnixNetDetect4(ctx AppContext) UnixNetDetect4 {
    return UnixNetDetect4{
        ctx: ctx,
    }
}


func (und UnixNetDetect4) detectNsHosts4(info *IPNetworkInfo) {
    var filepath = "/etc/resolv.conf"

    // Read the file
    var bytes, err = ioutil.ReadFile(filepath)
    if err != nil {
        und.ctx.ft.printError("Failed to detect ns servers", err)
        return
    }

    var content = string(bytes)
    und.parseResolvConf4(content, info)
}


func (und UnixNetDetect4) parseResolvConf4(content string,
                                            info *IPNetworkInfo) {

    var nameservers = und.parseResolvConf(content)

    for _, nameserver := range nameservers {
        var ip = net.ParseIP(nameserver)

        info.NsHosts = append(info.NsHosts, NsServer{
            Ip: ip,
        })
    }
}

func (und UnixNetDetect4) parseResolvConf(content string) []string {
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
