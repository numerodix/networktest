package main


type Pinger interface {
    ping(host string, cnt int, timeoutMs int) PingExecution
}
