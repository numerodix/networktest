package main

import "fmt"


func main() {
    var mgr = ProcMgr("ip", "addr")
    var res = mgr.run()
    fmt.Printf("%s\n", res)
}
