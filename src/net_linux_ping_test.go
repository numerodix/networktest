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


const linuxPing6Output = `
PING ::1(::1) 56 data bytes
64 bytes from ::1: icmp_seq=1 ttl=64 time=0.049 ms

--- ::1 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.049/0.049/0.049/0.000 ms
`


const bsdPing6Output = `
PING6(56=40+8+8 bytes) ::1 --> ::1
16 bytes from ::1, icmp_seq=0 hlim=64 time=0.433 ms

--- ::1 ping6 statistics ---
1 packets transmitted, 1 packets received, 0.0% packet loss
round-trip min/avg/max/std-dev = 0.433/0.433/0.433/0.000 ms
`


func Test_linuxParsePing4(t *testing.T) {
    var ctx = TestAppContext()
    var pinger = NewLinuxPinger(ctx)
    var pingExec = pinger.parsePing(linuxPing4Output)

    assertStrEq(t, "yahoo.com", pingExec.Host, "wrong host")
    assertStrEq(t, "154.327", strconv.FormatFloat(pingExec.Time, 'f', 3, 64), "wrong time")
    assertPtrEq(t, nil, pingExec.Err, "wrong err")
}


func Test_linuxParsePing6(t *testing.T) {
    var ctx = TestAppContext()
    var pinger = NewLinuxPinger(ctx)
    var pingExec = pinger.parsePing(linuxPing6Output)

    assertStrEq(t, "::1", pingExec.Host, "wrong host")
    assertStrEq(t, "0.049", strconv.FormatFloat(pingExec.Time, 'f', 3, 64), "wrong time")
    assertPtrEq(t, nil, pingExec.Err, "wrong err")
}


func Test_bsdParsePing6(t *testing.T) {
    var ctx = TestAppContext()
    var pinger = NewLinuxPinger(ctx)
    var pingExec = pinger.parsePing(bsdPing6Output)

    assertStrEq(t, "::1", pingExec.Host, "wrong host")
    assertStrEq(t, "0.433", strconv.FormatFloat(pingExec.Time, 'f', 3, 64), "wrong time")
    assertPtrEq(t, nil, pingExec.Err, "wrong err")
}
