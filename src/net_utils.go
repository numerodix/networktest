package main

//import "fmt"
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


func maskBytesToMask4(bytes []byte) net.IPMask {
    var mask = net.IPv4Mask(
        bytes[0],
        bytes[1],
        bytes[2],
        bytes[3])

    return mask
}
