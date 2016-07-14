/*
    Handles reading of /etc/resolv.conf
*/

package main

import (
    "fmt"
    "io/ioutil"
    "regexp"
    "strings"
)


func DetectNameservers() []string {
    var filepath = "/etc/resolv.conf"

    // Read the file
    var bytes, err = ioutil.ReadFile(filepath)
    if err != nil {
        fmt.Printf("Could not read %s: %s\n", filepath, err)
        return []string{}
    }

    // Parse the nameservers
    var nameservers []string
    var content = string(bytes)
    var lines = strings.Split(content, "\n")
    rx := regexp.MustCompile("nameserver ([^ ]*)")

    for i := range lines {
        var line = lines[i]

        if rx.MatchString(line) {
            var ns = rx.FindStringSubmatch(line)[1]
            nameservers = append(nameservers, ns)
        }
    }

    return nameservers
}
