/*
    Handles invocation of /sbin/ifconfig and parsing its output.
*/

package main

import (
    "bytes"
    "fmt"
    "net"
    "os"
    "os/exec"
    "regexp"
    "sort"
    "strings"
)


type IfaceBlock struct {
    Iface string
    LinkEncap string
    HWaddr string
    IPv4 string
    Broadcast string
    Mask string
    IPv6 string
    Scope string
    Status string
    Mtu string
}

// Sorting for []IfaceBlock
type ByIPv4 []IfaceBlock
func (ibs ByIPv4) Len() int {
    return len(ibs)
}
func (ibs ByIPv4) Swap(i, j int) {
    ibs[i], ibs[j] = ibs[j], ibs[i]
}
func (ibs ByIPv4) Less(i, j int) bool {
    var xIp = net.ParseIP(ibs[i].IPv4)
    var yIp = net.ParseIP(ibs[j].IPv4)
    return IPIsLesser(xIp, yIp)
}

type IfconfigExecution struct {
    IfaceBlocks []IfaceBlock
    Error error
}


func Ifconfig() IfconfigExecution {
    // Construct the args
    var executable = "/sbin/ifconfig"
    var args []string
    args = append(args, fmt.Sprintf("-a"))

    // Construct the cmd
    cmd := exec.Command(executable, args...)
    var out bytes.Buffer
    cmd.Stdout = &out

    // Invoke the cmd
    os.Setenv("LC_ALL", "C")
    err := cmd.Run()
    if err != nil {
        return IfconfigExecution{
            Error: fmt.Errorf("Failed to run ifconfig: %q", err),
        }
    }

    /* Output:
      $ /sbin/ifconfig -a
      docker0   Link encap:Ethernet  HWaddr 02:42:4d:ed:8b:26  
                inet addr:172.17.0.1  Bcast:0.0.0.0  Mask:255.255.0.0
                UP BROADCAST MULTICAST  MTU:1500  Metric:1
                RX packets:0 errors:0 dropped:0 overruns:0 frame:0
                TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
                collisions:0 txqueuelen:0 
                RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
      
      eth0      Link encap:Ethernet  HWaddr 14:da:e9:d5:3f:a2  
                inet addr:192.168.1.6  Bcast:192.168.1.255  Mask:255.255.255.0
                inet6 addr: fe80::16da:e9ff:fed5:3fa2/64 Scope:Link
                UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
                RX packets:19707243 errors:0 dropped:0 overruns:0 frame:0
                TX packets:14390808 errors:0 dropped:0 overruns:0 carrier:0
                collisions:0 txqueuelen:1000 
                RX bytes:25437088401 (25.4 GB)  TX bytes:1349325734 (1.3 GB)
      
      lo        Link encap:Local Loopback  
                inet addr:127.0.0.1  Mask:255.0.0.0
                inet6 addr: ::1/128 Scope:Host
                UP LOOPBACK RUNNING  MTU:65536  Metric:1
                RX packets:487316 errors:0 dropped:0 overruns:0 frame:0
                TX packets:487316 errors:0 dropped:0 overruns:0 carrier:0
                collisions:0 txqueuelen:0 
                RX bytes:93221522 (93.2 MB)  TX bytes:93221522 (93.2 MB)
      
      wlan0     Link encap:Ethernet  HWaddr 74:2f:68:ad:d6:23  
                inet addr:192.168.1.8  Bcast:192.168.1.255  Mask:255.255.255.0
                inet6 addr: fe80::762f:68ff:fead:d623/64 Scope:Link
                UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
                RX packets:18451 errors:0 dropped:0 overruns:0 frame:0
                TX packets:9246 errors:0 dropped:0 overruns:0 carrier:0
                collisions:0 txqueuelen:1000 
                RX bytes:3175200 (3.1 MB)  TX bytes:2228700 (2.2 MB)

    */

    // Parse the output into iface blocks
    var stdout = out.String()
    var blocks = strings.Split(stdout, "\n\n")
    var ifaceBlocks = []IfaceBlock{}
    rxIface := regexp.MustCompile("^[^ ]+")
    rxLinkEncap := regexp.MustCompile("Link encap:(.*)")
    rxHwAddr := regexp.MustCompile("HWaddr ([^ ]+)")
    rxIpv4 := regexp.MustCompile("inet addr:([^ ]+)")
    rxBcast := regexp.MustCompile("Bcast:([^ ]+)")
    rxMask := regexp.MustCompile("Mask:([^ ]+)")
    rxIpv6 := regexp.MustCompile("inet6 addr: ([^ ]+)")
    rxScope := regexp.MustCompile("Scope:([^ ]+)")
    rxPreMtu := regexp.MustCompile("(.*)MTU")
    rxMtu := regexp.MustCompile("MTU:([^ ]+)")

    for i := range blocks {
        var block = blocks[i]

        // If the block is empty skip it
        if strings.TrimSpace(block) == "" {
            continue
        }

        // Line starts with an iface
        var iface = rxIface.FindStringSubmatch(block)[0]

        var linkEncapPart = strings.Split(block[10:], "  ")[0]
        var linkEncap = rxLinkEncap.FindStringSubmatch(linkEncapPart)[1]

        var hwAddr = ""
        if rxHwAddr.MatchString(block) {
            hwAddr = rxHwAddr.FindStringSubmatch(block)[1]
        }

        var ipv4 = ""
        if rxIpv4.MatchString(block) {
            ipv4 = rxIpv4.FindStringSubmatch(block)[1]
        }

        var bcast = ""
        if rxBcast.MatchString(block) {
            bcast = rxBcast.FindStringSubmatch(block)[1]
        }

        var mask = ""
        if rxMask.MatchString(block) {
            mask = rxMask.FindStringSubmatch(block)[1]
        }

        var ipv6 = ""
        if rxIpv6.MatchString(block) {
            ipv6 = rxIpv6.FindStringSubmatch(block)[1]
        }

        var scope = ""
        if rxScope.MatchString(block) {
            scope = rxScope.FindStringSubmatch(block)[1]
        }

        var status = ""
        if rxPreMtu.MatchString(block) {
            status = rxPreMtu.FindStringSubmatch(block)[1]
        }

        var mtu = ""
        if rxMtu.MatchString(block) {
            mtu = rxMtu.FindStringSubmatch(block)[1]
        }

        // Exclude block if the ipv4 ip is not detected (we're not on the
        // network)
        if ipv4 == "" {
            continue
        }

        ifaceBlocks = append(ifaceBlocks, IfaceBlock{
            Iface: strings.TrimSpace(iface),
            LinkEncap: strings.TrimSpace(linkEncap),
            HWaddr: strings.TrimSpace(hwAddr),
            IPv4: strings.TrimSpace(ipv4),
            Broadcast: strings.TrimSpace(bcast),
            Mask: strings.TrimSpace(mask),
            IPv6: strings.TrimSpace(ipv6),
            Scope: strings.TrimSpace(scope),
            Status: strings.TrimSpace(status),
            Mtu: strings.TrimSpace(mtu),
        })
    }

    sort.Sort(ByIPv4(ifaceBlocks))

    return IfconfigExecution{
        IfaceBlocks: ifaceBlocks,
    }
}
