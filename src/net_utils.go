package main

import "net"


func ipIsLesser(x, y net.IP) bool {
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


func ipMaskToNet4(ip *net.IP, mask *net.IPMask) net.IPNet {
    var bytes = make([]byte, len(*ip))

    for i := range *ip {
        bytes[i] = (*ip)[i] & (*mask)[i]
    }

    var ipobj = net.IPv4(bytes[0], bytes[1], bytes[2], bytes[3])
    var ipnet = net.IPNet{IP: ipobj, Mask: *mask}

    return ipnet
}


func ipIPMaskToNet4(ip *net.IP, mask *net.IP) net.IPNet {
    var ipnet = net.IPNet{
        IP: *ip,
        Mask: net.IPv4Mask(
            (*mask)[12],
            (*mask)[13],
            (*mask)[14],
            (*mask)[15]),
    }

    return ipnet
}


func ipnetToMask4(ipnet *net.IPNet) net.IP {
    var mask = net.IPv4(
        ipnet.Mask[0],
        ipnet.Mask[1],
        ipnet.Mask[2],
        ipnet.Mask[3])

    return mask
}
