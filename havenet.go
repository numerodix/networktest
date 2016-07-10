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


func Ping(host string, cnt int, timeout int) (float64, error) {

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
        log.Fatal(err)
        return -1, err
    }

    // Parse the time= value
    var stdout = out.String()
    rx := regexp.MustCompile("time=([^ ]*)")
    var time_s = rx.FindStringSubmatch(stdout)[1]
    var time, err2 = strconv.ParseFloat(time_s, 64)
    if err2 != nil {
        log.Fatal(err2)
        return -1, err2
    }

    return time, nil
}


func main() {
//    hosts := []string{"yahoo.com", "google.com"}
    hosts := []string{"localhost"}

    for i := range hosts {
        time, err := Ping(hosts[i], 1, 2)
        if err != nil {
            log.Fatal(err)
        }

        log.Println(fmt.Sprintf("%-14s  ping: %.1f ms", hosts[i], time))
    }
}
