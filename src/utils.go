package main

import (
    "net"
)


func IPIsLesser(x, y net.IP) bool {
    // catch bad input
    if x == nil || y == nil {
        return false
    }

    // i actually ranges 0 -> 15
    // the last 4 bytes are the ipv4 address
    for i := range x {
        if x[i] < y[i] {
            return true
        }
    }

    return false
}


func NetworkContains(ip string, mask string, candidateIp string) bool {
    // catch bad input
    if ip == "" || mask == "" || candidateIp == "" {
        return false
    }

    var ipStruct = net.ParseIP(ip)

    var maskIP = net.ParseIP(mask)
    var maskStruct = net.IPv4Mask(maskIP[12], maskIP[13], maskIP[14], maskIP[15])

    var netStruct = net.IPNet{IP:ipStruct, Mask:maskStruct}

    var candidateStruct = net.ParseIP(candidateIp)

    return netStruct.Contains(candidateStruct)
}
