package main

import "flag"
import "fmt"


func main() {
    // Parse command line args
    var _ = flag.Bool("4", true, "Test IPv4 network connectivity")
    var flagIpv6Ptr = flag.Bool("6", false, "Test IPv6 network connectivity")
    var noColorPtr = flag.Bool("nc", false, "No color output")
    var versionPtr = flag.Bool("V", false, "Display the version")
    flag.Parse()

    // Only print the version and exit
    if *versionPtr {
        fmt.Printf("%s\n", appVersion)
        return
    }

    var ipver = 4
    if *flagIpv6Ptr {
        ipver = 6
    }

    var ui = NetworkDetector(ipver, *noColorPtr)
    ui.run()
}
