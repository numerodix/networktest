package main

import (
    "bytes"
//    "errors"
    "fmt"
//    "log"
    "os/exec"
    "regexp"
    "strconv"
//    "strings"
)


type CommandResult struct {
    Id string
    FValue float64
    SValue string
    Error error
}


func Ping(ch chan CommandResult, host string, cnt int, timeout int) {
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
        ch <- CommandResult{Id: host, Error: err}
        return
    }

    // Parse the time= value
    var stdout = out.String()
    rx := regexp.MustCompile("time=([^ ]*)")
    var time_s = rx.FindStringSubmatch(stdout)[1]
    var time, err2 = strconv.ParseFloat(time_s, 64)
    if err2 != nil {
        ch <- CommandResult{Id: host, Error: err2}
        return
    }

    ch <- CommandResult{Id: host, FValue: time}
}


func main() {
//    hosts := []string{"yahoo.com", "google.com"}
    hosts := []string{
        "192.168.1.1",
        "192.228.79.201",
        "127.0.0.1",
        "127.0.1.1",
        "localhost",
        "yahoo.com",
        "google.com",
        "juventuz.com",
        "twitter.com",
        "facebook.com",
        "gmail.com",
        "golang.org",
        "www.nu.nl",
        "www.aftenposten.no",
        "www.bonjourchine.com",
        "github.com",
        "youtube.com",
    }
//    hosts := []string{"localhost"}
    ch := make(chan CommandResult)

    // Launch
    for i := range hosts {
        go Ping(ch, hosts[i], 1, 2)
    }

    // Collect
    sum := 0.0
    for i := range hosts {
        cmdres := <-ch

        if cmdres.Error != nil {
            fmt.Printf("Err: %s: %s\n", cmdres.Id, cmdres.Error)
            continue
        }

        sum += cmdres.FValue
        fmt.Printf("%-2d  %-34s: %.1f ms\n", i, cmdres.Id, cmdres.FValue)
    }

    fmt.Printf("Total time: %.1f ms\n", sum)
}
