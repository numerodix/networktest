package main

import (
    "bytes"
//    "errors"
    "fmt"
    "log"
    "os/exec"
    "regexp"
    "strconv"
//    "strings"
)


type CommandResult struct {
    FValue float64
    SValue string
    Error error
}


func Ping(pingch chan CommandResult, host string, cnt int, timeout int) {
    // Prepare return variables
    etime := -1.0

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
    err := cmd.Run()
    if err != nil {
        pingch <- CommandResult{FValue: etime, Error: err}
        return
    }

    // Parse the time= value
    var stdout = out.String()
    rx := regexp.MustCompile("time=([^ ]*)")
    var time_s = rx.FindStringSubmatch(stdout)[1]
    var time, err2 = strconv.ParseFloat(time_s, 64)
    if err2 != nil {
        pingch <- CommandResult{FValue: etime, Error: err2}
        return
    }

    pingch <- CommandResult{FValue: time, Error: nil}
}


func main() {
//    hosts := []string{"yahoo.com", "google.com"}
    hosts := []string{"localhost"}

    for i := range hosts {
        pingch := make(chan CommandResult)
        go Ping(pingch, hosts[i], 1, 2)
        pingres := <-pingch

        log.Println("time: ", pingres.FValue)
    }
}
