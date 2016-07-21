package main

import "fmt"
import "runtime"


func main() {
    if runtime.GOOS == "linux" {
        var info = linuxDetectNetConn4()
        fmt.Printf("%s\n", info)

    } else {
        fmt.Printf("Platform not supported: %s\n", runtime.GOOS)
    }
}
