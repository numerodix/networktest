package main

//import "fmt"
import "net"
import "sort"


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

// Sorting for Network
type ByNetwork []Network
func (nets ByNetwork) Len() int {
    return len(nets)
}
func (nets ByNetwork) Swap(i, j int) {
    nets[i], nets[j] = nets[j], nets[i]
}
func (nets ByNetwork) Less(i, j int) bool {
    return ipIsLesser(nets[i].Ip.IP, nets[j].Ip.IP)
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

func (ipa *IpAddr) getAsIpnet() net.IPNet {
    var ipnet = ipIPMaskToNet4(&ipa.Ip, &ipa.Mask)
    return ipnet
}

func (ipa *IpAddr) ipAsString() string {
    return ipa.Ip.String()
}

func (ipa *IpAddr) maskAsString() string {
    return ipa.Mask.String()
}

// Sorting for IpAddr
type ByIpAddr []IpAddr
func (ipas ByIpAddr) Len() int {
    return len(ipas)
}
func (ipas ByIpAddr) Swap(i, j int) {
    ipas[i], ipas[j] = ipas[j], ipas[i]
}
func (ipas ByIpAddr) Less(i, j int) bool {
    return ipIsLesser(ipas[i].Ip, ipas[j].Ip)
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

// Sorting for Gateway
type ByGateway []Gateway
func (gws ByGateway) Len() int {
    return len(gws)
}
func (gws ByGateway) Swap(i, j int) {
    gws[i], gws[j] = gws[j], gws[i]
}
func (gws ByGateway) Less(i, j int) bool {
    return ipIsLesser(gws[i].Ip, gws[j].Ip)
}


type NsServer struct {
    Ip net.IP

    // Is4
    // Is6

    // IpAsString
}

func (nshost *NsServer) ipAsString() string {
    return nshost.Ip.String()
}

// Sorting for NsServer
type ByNsServer []NsServer
func (nss ByNsServer) Len() int {
    return len(nss)
}
func (nss ByNsServer) Swap(i, j int) {
    nss[i], nss[j] = nss[j], nss[i]
}
func (nss ByNsServer) Less(i, j int) bool {
    return ipIsLesser(nss[i].Ip, nss[j].Ip)
}


type IP4NetworkInfo struct {
    Nets []Network
    Ips []IpAddr
    Gws []Gateway
    NsHosts []NsServer
    Errs []error

    // GetAllIfaces()
    // GetIpsForGw(gw)
    // GetNetsForIface(iface)
    // 
}

func (info *IP4NetworkInfo) getSortedNets() []Network {
    var nets = info.Nets
    sort.Sort(ByNetwork(nets))
    return nets
}

func (info *IP4NetworkInfo) getSortedIps() []IpAddr {
    var ips = info.Ips
    sort.Sort(ByIpAddr(ips))
    return ips
}

func (info *IP4NetworkInfo) getSortedGws() []Gateway {
    var gws = info.Gws
    sort.Sort(ByGateway(gws))
    return gws
}

func (info *IP4NetworkInfo) getSortedNsHosts() []NsServer {
    var nss = info.NsHosts
    sort.Sort(ByNsServer(nss))
    return nss
}


func (info *IP4NetworkInfo) getIpsForGw(gw *Gateway) []IpAddr {
    var ips = []IpAddr{}

    for _, ip := range info.Ips {
        var ipnet = ip.getAsIpnet()
        if ipnet.Contains(gw.Ip) {
            ips = append(ips, ip)
        }
    }

    return ips
}

func (info *IP4NetworkInfo) haveLocalNet() bool {
    // If we have a gateway and an ip on that network we have a local network
    // connection
    for _, gw := range info.getSortedGws() {
        var ips = info.getIpsForGw(&gw)
        if len(ips) != 0 {
            return true
        }
    }

    return false
}

func (info *IP4NetworkInfo) normalize() {
    var gwIps = make(map[string]int)
    var nsIps = make(map[string]int)
    var gws = []Gateway{}
    var nss = []NsServer{}

    // Remove duplicate gateways
    for _, gw := range info.Gws {
        var key = gw.Ip.String()
        gwIps[key] += 1

        // If this gateway already exists once, skip it
        if gwIps[key] <= 1 {
            gws = append(gws, gw)
        }
    }

    // Remove duplicate nshosts
    for _, nshost := range info.NsHosts {
        var key = nshost.Ip.String()
        nsIps[key] += 1

        // If this nshost already exists once, skip it
        if nsIps[key] <= 1 {
            nss = append(nss, nshost)
        }
    }

    info.Gws = gws
    info.NsHosts = nss
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
