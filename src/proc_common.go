package main


type ProcessResult struct {
    exitCode int
    stdout string
    stderr string
    err error
}
