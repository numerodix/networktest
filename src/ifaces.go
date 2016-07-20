package main

import "net"


type Interface struct {
    Name string
}


type Network struct {
    Iface Interface
    Ip net.IPNet

    // Is4
    // Is6

    // IpAsString
    // MaskAsString
    // BcastAsString
}


type IpAddr struct {
    Iface Interface
    Ip net.IP
    Mask net.IP

    // Is4
    // Is6

    // IpAsString
    // MaskAsString
}


type Gateway struct {
    Iface Interface
    Ip net.IP
    // Mask?

    // Is4
    // Is6

    // IpAsString
}


type NsServer struct {
    Ip net.IP

    // Is4
    // Is6

    // IpAsString
}


type IP4NetworkInfo struct {
    Nets []Network
    Ips []IpAddr
    Gws []Gateway
    NsHosts []NsServer

    // GetAllIfaces()
    // GetNetsForIface(iface)
    // 
}


type IP6NetworkInfo struct {
    Nets []Network
    Ips []IpAddr
    Gws []Gateway
    NsHosts []NsServer

    // GetAllIfaces()
    // GetNetsForIface(iface)
    // 
}
