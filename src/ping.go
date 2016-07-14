/*
    Handles invocation of ping and parsing its output.
*/

package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "regexp"
    "strconv"
)


type PingExecution struct {
    Host string
    Time float64
    Error error
}


func Ping(host string, cnt int, timeout int) PingExecution {
    // Construct the args
    var executable = "ping"
    var args []string
    args = append(args, fmt.Sprintf("-c%d", cnt))
    args = append(args, fmt.Sprintf("-W%d", timeout))
    args = append(args, host)

    // Construct the cmd
    cmd := exec.Command(executable, args...)
    var out bytes.Buffer
    cmd.Stdout = &out

    // Invoke the cmd
    os.Setenv("LC_ALL", "C")
    err := cmd.Run()
    if err != nil {
        return PingExecution{
            Host: host,
            Error: fmt.Errorf("Failed to run ping: %q", err),
        }
    }

    /* Output:
      $ ping -c1 -W2 localhost
      PING localhost (127.0.0.1) 56(84) bytes of data.
      64 bytes from localhost (127.0.0.1): icmp_seq=1 ttl=64 time=0.061 ms

      --- localhost ping statistics ---
      1 packets transmitted, 1 received, 0% packet loss, time 0ms
      rtt min/avg/max/mdev = 0.061/0.061/0.061/0.000 ms
    */

    // Parse the time= value
    var stdout = out.String()
    rx := regexp.MustCompile("time=([^ ]*)")
    var time_s = rx.FindStringSubmatch(stdout)[1]
    var time, err2 = strconv.ParseFloat(time_s, 64)
    if err2 != nil {
        return PingExecution{
            Host: host,
            Error: fmt.Errorf("Failed to parse ping output: %q", err2),
        }
    }

    // Return success
    return PingExecution{Host: host, Time: time}
}
