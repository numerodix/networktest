package main

import "io/ioutil"
import "net"
import "regexp"
import "strings"


func unixDetectNsHosts4(info *IP4NetworkInfo) {
    var filepath = "/etc/resolv.conf"

    // Read the file
    var bytes, err = ioutil.ReadFile(filepath)
    if err != nil {
        // XXX print some kind of useful error
        return
    }

    var content = string(bytes)
    unixParseResolvConf4(content, info)
}


func unixParseResolvConf4(content string, info *IP4NetworkInfo) {
    var nameservers = unixParseResolvConf(content)

    for i := range nameservers {
        var nameserver = nameservers[i]
        var ip = net.ParseIP(nameserver)

        info.NsHosts = append(info.NsHosts, NsServer{
            Ip: ip,
        })
    }
}

func unixParseResolvConf(content string) []string {
    // Parse the nameservers
    var nameservers []string
    var lines = strings.Split(content, "\n")
    rx := regexp.MustCompile("nameserver ([^ ]*)")

    for i := range lines {
        var line = lines[i]

        if rx.MatchString(line) {
            var ns = rx.FindStringSubmatch(line)[1]
            nameservers = append(nameservers, ns)
        }
    }

    return nameservers
}
