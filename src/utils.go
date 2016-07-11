package main

import (
    "net"
)


func LessIPs(x, y net.IP) bool {
    // i actually ranges 0 -> 15
    // the last 4 bytes are the ipv4 address
    for i := range x {
        if x[i] < y[i] {
            return true
        }
    }

    return false
}
