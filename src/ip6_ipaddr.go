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


type Ip6IpAddrBlock struct {
    Iface string
    IPv6 net.IP
    Network net.IPNet
    Scope string
}

type Ip6IpAddrExecution struct {
    Ip6IpAddrBlocks []Ip6IpAddrBlock
    Error error
}


func Ip6IpAddr() Ip6IpAddrExecution {
    // Construct the args
    var executable = "/sbin/ip"
    var args []string
    args = append(args, fmt.Sprintf("-6"))
    args = append(args, fmt.Sprintf("addr"))
    args = append(args, fmt.Sprintf("show"))

    // Construct the cmd
    cmd := exec.Command(executable, args...)
    var out bytes.Buffer
    cmd.Stdout = &out

    // Invoke the cmd
    os.Setenv("LC_ALL", "C")
    err := cmd.Run()
    if err != nil {
        return Ip6IpAddrExecution{
            Error: fmt.Errorf("Failed to run %s: %q", executable, err),
        }
    }

    /* Output:
      $ /sbin/ip -6 addr show
      1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 
          inet6 ::1/128 scope host 
             valid_lft forever preferred_lft forever
      2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qlen 1000
          inet6 2a00:dd80:da::e27/64 scope global 
             valid_lft forever preferred_lft forever
          inet6 fe80::16da:fae1:c9ea:a4b9/64 scope link 
             valid_lft forever preferred_lft forever
      3: wlan0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qlen 1000
          inet6 fe80::762f:fe64:b7c7:7b7a/64 scope link 
             valid_lft forever preferred_lft forever
    */

    // Parse the output into lines
    var stdout = out.String()
    var lines = strings.Split(stdout, "\n")

    rxIface := regexp.MustCompile("^[0-9]+: ([^ ]+):")
    rxInet6 := regexp.MustCompile(
        "^[ ]{4}inet6 ([a-f0-9:]+)/([0-9]+) scope ([A-Za-z0-9]+)")

    var ip6IpAddrBlocks = []Ip6IpAddrBlock{}

    // loop variables
    var iface = ""
    var ipv6 = ""
    var mask = ""
    var scope = ""

    for i := range lines {
        var line = lines[i]

        if rxIface.MatchString(line) {
            iface = rxIface.FindStringSubmatch(line)[1]
        }

        if rxInet6.MatchString(line) {
            ipv6 = rxInet6.FindStringSubmatch(line)[1]
            mask = rxInet6.FindStringSubmatch(line)[2]
            scope = rxInet6.FindStringSubmatch(line)[3]

            var ipNet = fmt.Sprintf("%s/%s", ipv6, mask)
            var ip = net.ParseIP(ipv6)
            var _, ipnet, err = net.ParseCIDR(ipNet)

            if err != nil {
                return Ip6IpAddrExecution{
                    Error: fmt.Errorf("Failed to parse ipnet %s: %q", ipNet, err),
                }
            }

            ip6IpAddrBlocks = append(ip6IpAddrBlocks, Ip6IpAddrBlock{
                Iface: iface,
                IPv6: ip,
                Network: *ipnet,
                Scope: scope,
            })
        }
    }

    return Ip6IpAddrExecution{
        Ip6IpAddrBlocks: ip6IpAddrBlocks,
    }
}
