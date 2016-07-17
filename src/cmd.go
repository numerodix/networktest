package main

import (
    "flag"
)


func main() {
    // Parse command line args
    flagIpv4Ptr := flag.Bool("4", true, "Test IPv4 network connectivity")
    flagIpv6Ptr := flag.Bool("6", false, "Test IPv6 network connectivity")
    flag.Parse()

    // Decide on execution params
    var ipv4 = *flagIpv4Ptr
    if *flagIpv6Ptr {
        ipv4 = !*flagIpv6Ptr
    }

    // Run the program
    if ipv4 {
        HaveNet4()
    } else {
        HaveNet6()
    }
}
