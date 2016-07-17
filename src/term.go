package main

import (
    "os"
)


/*
    Detect whether we are connected to a terminal (TERM set) and whether the
    terminal is dumb (does not support ansi control codes).
*/
func TerminalIsDumb() bool {
    var term = os.Getenv("TERM")

    if term == "" || term == "dumb" {
        return true
    }

    return false
}
