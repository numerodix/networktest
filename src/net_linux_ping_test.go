package main

//import "fmt"
//import "strings"
import "strconv"
import "testing"


const linuxPing4Output = `
PING yahoo.com (98.138.253.109) 56(84) bytes of data.
64 bytes from ir1.fp.vip.ne1.yahoo.com (98.138.253.109): icmp_seq=1 ttl=48 time=154 ms

--- yahoo.com ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 154.327/154.327/154.327/0.000 ms
`


func Test_linuxParsePing4(t *testing.T) {
    var ctx = TestAppContext()
    var pinger = NewLinuxPinger(ctx)
    var pingExec = pinger.parsePing(linuxPing4Output)

    assertStrEq(t, "yahoo.com", pingExec.Host, "wrong host")
    assertStrEq(t, "154.327", strconv.FormatFloat(pingExec.Time, 'f', 3, 64), "wrong time")
    assertPtrEq(t, nil, pingExec.Err, "wrong err")
}
