package main

//import "fmt"
import "os"
import "path"
import "strings"


type Toolbox struct {
    osName string
}


func (tb *Toolbox) haveTool(command string) bool {
    var envPath = os.Getenv("PATH")

    var delim = ":"
    if tb.osName == "windows" {
        delim = ";"
    }

    var dirpaths = strings.Split(envPath, delim)
    for _, dirpath := range dirpaths {
        var thisPath = path.Join(dirpath, command)

        // If the file does not exist stat() will fail
        var _, err = os.Stat(thisPath)
        if !os.IsNotExist(err) {
            return true
        }
    }

    return false
}
