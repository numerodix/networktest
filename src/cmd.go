package main

import "flag"


func main() {
    // Parse command line args
    var _ = flag.Bool("4", true, "Test IPv4 network connectivity")
    var flagIpv6Ptr = flag.Bool("6", false, "Test IPv6 network connectivity")
    flag.Parse()

    var ipver = 4
    if *flagIpv6Ptr {
        ipver = 6
    }

    var ui = NetworkDetector(ipver)
    run(&ui)
}
