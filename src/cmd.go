package main

import "fmt"


func main() {
    var mgr = ProcMgr("ip", "addr")
    var res = mgr.run()

    fmt.Printf("%s\n", res.exitCode)
    fmt.Printf("%s\n", res.stderr)
    fmt.Printf("%s\n", res.err)
}
