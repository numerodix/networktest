package main

import (
    "net"
)


func LessIPs(x, y net.IP) bool {
    // n actually ranges 0 -> 15
    // the last 4 bytes are the ipv4 address
    for n := range x {
        if x[n] < y[n] {
            return true
        }
    }

    return false
}
