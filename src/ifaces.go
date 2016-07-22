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

func (net *Network) ipAsString() string {
    return net.Ip.IP.String()
}

func (net *Network) maskAsString() string {
    var mask = ipnetToMask4(&net.Ip)
    return mask.String()
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

func (ipa *IpAddr) ipAsString() string {
    return ipa.Ip.String()
}

func (ipa *IpAddr) maskAsString() string {
    return ipa.Mask.String()
}


type Gateway struct {
    Iface Interface
    Ip net.IP
    // Mask?

    // Is4
    // Is6

    // IpAsString
}

func (gw *Gateway) ipAsString() string {
    return gw.Ip.String()
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
    Errs []error

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
