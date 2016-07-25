package main

//import "fmt"
//import "strings"
import "strconv"
import "testing"


const winPing4Output = `
Pinging yahoo.com [206.190.36.45] with 32 bytes of data:
Reply from 206.190.36.45: bytes=32 time=193ms TTL=51

Ping statistics for 206.190.36.45:
    Packets: Sent = 1, Received = 1, Lost = 0 (0% loss),
Approximate round trip times in milli-seconds:
    Minimum = 193ms, Maximum = 193ms, Average = 193ms
`


func Test_winParsePing4(t *testing.T) {
    var ctx = TestAppContext()
    var pinger = NewWinPinger4(ctx)
    var pingExec = pinger.parsePing4(winPing4Output)

    assertStrEq(t, "yahoo.com", pingExec.Host, "wrong host")
    assertStrEq(t, "193", strconv.FormatFloat(pingExec.Time, 'f', 0, 64), "wrong time")
    assertPtrEq(t, nil, pingExec.Err, "wrong err")
}
