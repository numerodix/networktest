package main

import "fmt"
import "net"


func ipIs4(ip net.IP) bool {
    // Guard against bad input
    if ip == nil {
        panic("Input cannot be nil")
    }

    if ip.To4() == nil {
        return false
    }

    return true
}

func ipIs6(ip net.IP) bool {
    // Guard against bad input
    if ip == nil {
        panic("Input cannot be nil")
    }

    // To16 will not return nil the way that To4 will. So we need to negate the
    // test for ipv4 instead.
    if ipIs4(ip) {
        return false
    }

    return true
}


func maskIs4(mask net.IPMask) bool {
    // Guard against bad input
    if mask == nil {
        panic("Input cannot be nil")
    }

    if len(mask) == 4 {
        return true
    }

    return false
}

func maskIs6(mask net.IPMask) bool {
    // Guard against bad input
    if mask == nil {
        panic("Input cannot be nil")
    }

    if len(mask) == 16 {
        return true
    }

    return false
}


func ipIsLesser(x, y net.IP) bool {
    // catch bad input
    if x == nil || y == nil {
        panic("Inputs cannot be nil")
    }

    // i actually ranges 0 -> 15
    // the last 4 bytes are the ipv4 address
    for i := range x {
        // if the current (high order) byte is lesser the whole ip is lesser
        if x[i] < y[i] {
            return true
        }

        // if the current (high order) byte is greater the whole ip cannot be
        // lesser
        if x[i] > y[i] {
            break
        }

        // otherwise we have a tie and we keep looping
    }

    return false
}


func applyMask(ip *net.IP, mask *net.IPMask) net.IPNet {
    // catch bad input
    if ip == nil || mask == nil {
        panic("Inputs cannot be nil")
    }

    // Make sure the inputs are the same ip version
    if ipIs4(*ip) && maskIs6(*mask) {
        panic("Inputs do not match")
    }
    if ipIs6(*ip) && maskIs4(*mask) {
        panic("Inputs do not match")
    }

    var ipobj = ip.Mask(*mask)
    var ipnet = net.IPNet{IP: ipobj, Mask: *mask}

    return ipnet
}


func maskAsIpToIPMask(mask *net.IP) net.IPMask {
    var length int
    var offset int

    if ipIs4(*mask) {
        length = 4
        offset = 12  // skip the first 12 bytes of the array
    } else {
        length = 16
        offset = 0
    }

    // Count the bits in the netmask
    var bits = 0
    for i := 0; i < length; i++ {
        var octet = (*mask)[i + offset]
        bits += (int(octet) + 1) >> 5
    }

    var ipmask = net.CIDRMask(bits, length * 8)

    return ipmask
}


func ipAndMaskToIPNet(ip *net.IP, mask *net.IP) net.IPNet {
    // catch bad input
    if ip == nil || mask == nil {
        panic("Inputs cannot be nil")
    }

    // Make sure the inputs are the same ip version
    if ipIs4(*ip) && ipIs6(*mask) {
        panic("Inputs do not match")
    }
    if ipIs6(*ip) && ipIs4(*mask) {
        panic("Inputs do not match")
    }

    var ipmask = maskAsIpToIPMask(mask)

    var ipnet = net.IPNet{
        IP: *ip,
        Mask: ipmask,
    }

    return ipnet
}


func ipnetMaskAsIP(ipnet *net.IPNet) net.IP {
    var ipobj net.IP

    if maskIs4(ipnet.Mask) {
        ipobj = net.IPv4(
                    ipnet.Mask[0],
                    ipnet.Mask[1],
                    ipnet.Mask[2],
                    ipnet.Mask[3])

    } else {
        var ipStr = fmt.Sprintf("%x%x:%x%x:%x%x:%x%x:%x%x:%x%x:%x%x:%x%x",
                            ipnet.Mask[0],
                            ipnet.Mask[1],
                            ipnet.Mask[2],
                            ipnet.Mask[3],
                            ipnet.Mask[4],
                            ipnet.Mask[5],
                            ipnet.Mask[6],
                            ipnet.Mask[7],
                            ipnet.Mask[8],
                            ipnet.Mask[9],
                            ipnet.Mask[10],
                            ipnet.Mask[11],
                            ipnet.Mask[12],
                            ipnet.Mask[13],
                            ipnet.Mask[14],
                            ipnet.Mask[15])
        ipobj = net.ParseIP(ipStr)
    }

    return ipobj
}


func maskBytesToIPMask(bytes []byte) net.IPMask {
    var mask = net.IPv4Mask(
        bytes[0],
        bytes[1],
        bytes[2],
        bytes[3])

    return mask
}
