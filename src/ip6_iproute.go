package main

import (
    "bytes"
    "fmt"
    "net"
    "os"
    "os/exec"
    "regexp"
    "strings"
)


type Ip6RouteBlock struct {
    Iface string
    IPv6 net.IP
}

type Ip6RouteExecution struct {
    Ip6RouteBlocks []Ip6RouteBlock
    Error error
}


func Ip6Route() Ip6RouteExecution {
    // Construct the args
    var executable = "/sbin/ip"
    var args []string
    args = append(args, fmt.Sprintf("-6"))
    args = append(args, fmt.Sprintf("route"))
    args = append(args, fmt.Sprintf("show"))

    // Construct the cmd
    cmd := exec.Command(executable, args...)
    var out bytes.Buffer
    cmd.Stdout = &out

    // Invoke the cmd
    os.Setenv("LC_ALL", "C")
    err := cmd.Run()
    if err != nil {
        return Ip6RouteExecution{
            Error: fmt.Errorf("Failed to run %s: %q", executable, err),
        }
    }

    /* Output:
      $ ip -6 route show 
      2a00:ab41:e1::/64 dev eth0  proto kernel  metric 256 
      fe80::/64 dev eth0  proto kernel  metric 256 
      default via 2a00:ab41:e1::1 dev eth0  metric 1024 
    */

    // Parse the output into lines
    var stdout = out.String()
    var lines = strings.Split(stdout, "\n")

    rxIface := regexp.MustCompile("^default via ([A-Fa-f0-9:]+) dev ([^ ]+)")

    var ip6RouteBlocks = []Ip6RouteBlock{}

    // loop variables
    var iface = ""
    var ipv6 = ""

    for i := range lines {
        var line = lines[i]

        if rxIface.MatchString(line) {
            ipv6 = rxIface.FindStringSubmatch(line)[1]
            iface = rxIface.FindStringSubmatch(line)[2]

            var ip = net.ParseIP(ipv6)

            ip6RouteBlocks = append(ip6RouteBlocks, Ip6RouteBlock{
                Iface: iface,
                IPv6: ip,
            })
        }
    }

    return Ip6RouteExecution{
        Ip6RouteBlocks: ip6RouteBlocks,
    }
}
