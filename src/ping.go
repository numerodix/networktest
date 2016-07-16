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
    "time"
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

    // Launch the cmd
    os.Setenv("LC_ALL", "C")
    err := cmd.Start()
    if err != nil {
        return PingExecution{
            Host: host,
            Error: fmt.Errorf("Failed to run ping: %q", err),
        }
    }

    // Poll for it to complete:
    // We can't really count on ping's -W parameter to bound the execution of
    // ping itself, because -W limits how long ping will wait for the response
    // packet. But before the ICMP packet can be sent the hostname has to be
    // resolved using DNS, and that's not included in the bound set by -W. So
    // instead we launch the process and kill the ping process once timeout has
    // been reached. This *should* guarantee that this function always honors
    // the intended timeout.
    var elapsedMs = 0
    for {
        // Use the pid to detect if the process exists
        var _, err = os.FindProcess(cmd.Process.Pid)

        // If it doesn't:
        // - it's still forking (not yet running)
        // - or it exited already (we check stdout to see if it's non-empty)
        if err == nil && out.String() != "" {
            break
        }

        // We reached the timeout, so kill the process and return an error
        if elapsedMs >= (timeout * 1000) {
            err = cmd.Process.Kill()
            return PingExecution{
                Host: host,
                Error: fmt.Errorf("Failed to run ping: %q", err),
            }
        }

        // The process hasn't finished yet, so wait a while and loop around
        time.Sleep(50 * time.Millisecond)
        elapsedMs += 50
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
    if !rx.MatchString(stdout) {
        return PingExecution{
            Host: host,
            Error: fmt.Errorf("Failed to parse ping output: %q", stdout),
        }
    }

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
