/*
    Handles invocation of /sbin/route and parsing its output.
*/

package main

import (
    "bytes"
//    "errors"
    "fmt"
//    "log"
    "net"
    "os"
    "os/exec"
    "regexp"
//    "strconv"
    "sort"
    "strings"
)


type RouteLine struct {
    Destination string
    Gateway string
    Genmask string
    Flags string
    Metric string
    Ref string
    Use string
    Iface string
}

type RouteExecution struct {
    Lines []RouteLine
    Error error
}


type Network struct {
    Iface string
    Network string
    Netmask string
}

// Sorting for []Network
type ByNetwork []Network
func (ns ByNetwork) Len() int {
    return len(ns)
}
func (ns ByNetwork) Swap(i, j int) {
    ns[i], ns[j] = ns[j], ns[i]
}
func (ns ByNetwork) Less(i, j int) bool {
    var xIp = net.ParseIP(ns[i].Network)
    var yIp = net.ParseIP(ns[j].Network)
    return LessIPs(xIp, yIp)
}

type Gateway struct {
    Iface string
    Gateway string
}


func Route() RouteExecution {
    // Construct the args
    var executable = "/sbin/route"
    var args []string
    args = append(args, fmt.Sprintf("-n"))

    // Construct the cmd
    cmd := exec.Command(executable, args...)
    var out bytes.Buffer
    cmd.Stdout = &out

    // Invoke the cmd
    os.Setenv("LC_ALL", "C")
    err := cmd.Run()
    if err != nil {
        return RouteExecution{
            Error: fmt.Errorf("Failed to run route: %q", err),
        }
    }

    /* Output:
      $ /sbin/route -n
      Kernel IP routing table
      Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
      0.0.0.0         192.168.1.1     0.0.0.0         UG    0      0        0 eth0
      172.17.0.0      0.0.0.0         255.255.0.0     U     0      0        0 docker0
      192.168.1.0     0.0.0.0         255.255.255.0   U     1      0        0 eth0
      192.168.1.0     0.0.0.0         255.255.255.0   U     9      0        0 wlan0
    */

    // Parse the output into lines
    var stdout = out.String()
    var lines = strings.Split(stdout, "\n")
    var routeLines = []RouteLine{}

    for i := range lines {
        var line = lines[i]

        // Skip the first two lines
        if i < 2 {
            continue
        }

        // If the line is empty skip it
        if strings.TrimSpace(line) == "" {
            continue
        }

        // Pick out each field using a constant width
        var destination = line[:16]
        var gateway = line[16:32]
        var genmask = line[32:48]
        var flags = line[48:54]
        var metric = line[54:61]
        var ref = line[61:68]
        var use = line[68:72]
        var iface = line[72:]

        routeLines = append(routeLines, RouteLine{
            Destination: strings.TrimSpace(destination),
            Gateway: strings.TrimSpace(gateway),
            Genmask: strings.TrimSpace(genmask),
            Flags: strings.TrimSpace(flags),
            Metric: strings.TrimSpace(metric),
            Ref: strings.TrimSpace(ref),
            Use: strings.TrimSpace(use),
            Iface: strings.TrimSpace(iface),
        })
    }

    return RouteExecution{
        Lines: routeLines,
    }
}


func (routeExec RouteExecution) GetNetworks() []Network {
    rx := regexp.MustCompile("^[1-9]")
    var networks = []Network{}

    for i := range routeExec.Lines {
        var line = routeExec.Lines[i]

        if rx.MatchString(line.Destination) {
            networks = append(networks, Network{
                Iface: line.Iface,
                Network: line.Destination,
                Netmask: line.Genmask,
            })
        }
    }

    sort.Sort(ByNetwork(networks))

    return networks
}

func (routeExec RouteExecution) GetGateways() []Gateway {
    rx := regexp.MustCompile("UG")
    var gateways = []Gateway{}

    for i := range routeExec.Lines {
        var line = routeExec.Lines[i]

        if rx.MatchString(line.Flags) {
            gateways = append(gateways, Gateway{
                Iface: line.Iface,
                Gateway: line.Gateway,
            })
        }
    }

    return gateways
}
