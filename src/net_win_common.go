package main

import "fmt"
import "regexp"


// Assign names like eth0, wlan0 based on the section name
type IfaceNamer struct {
    dict map[string]int
}

func InterfaceNamer() IfaceNamer {
    return IfaceNamer{
        dict: make(map[string]int),
    }
}

func (in *IfaceNamer) getPrefix(section string) string {
    var rxEth = regexp.MustCompile("(?i)ethernet")
    var rxWifi = regexp.MustCompile("(?i)wireless")

    if rxWifi.MatchString(section) {
        return "wlan"
    }
    if rxEth.MatchString(section) {
        return "eth"
    }

    return "if"
}

func (in *IfaceNamer) allocateName(section string) string {
    var prefix = in.getPrefix(section)

    var cnt = in.dict[prefix]
    cnt += 1
    in.dict[prefix] = cnt

    var name = fmt.Sprintf("%s%d", prefix, cnt)

    return name
}
