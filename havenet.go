package main

import (
//    "bytes"
    "errors"
//    "fmt"
    "log"
//    "os/exec"
//    "regexp"
//    "strconv"
//    "strings"
)


type CommandResult struct {
    FValue float64
    SValue string
    Error error
}


func Ping(pingch chan CommandResult) {
    fvalue := 0.0
    err := errors.New("")

    pingch <- CommandResult{FValue: fvalue, Error: err}
}


func main() {
//    hosts := []string{"yahoo.com", "google.com"}
//    hosts := []string{"localhost"}

    pingch := make(chan CommandResult)
    go Ping(pingch)
    pingres := <-pingch

    log.Println("time: ", pingres.FValue)
}
