package main

import "net"


func ipMaskToNet4(ip net.IP, mask net.IPMask) net.IPNet {
    var bytes = make([]byte, len(ip))

    for i := range ip {
        bytes[i] = ip[i] & mask[i]
    }

    var ipobj = net.IPv4(bytes[0], bytes[1], bytes[2], bytes[3])
    var ipnet = net.IPNet{IP: ipobj, Mask: mask}

    return ipnet
}
