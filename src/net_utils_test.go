package main

//import "fmt"
import "net"
import "testing"


func Test_ipIs4(t *testing.T) {
    var ip4 = net.ParseIP("1.1.1.1")
    var ip6 = net.ParseIP("::1")

    assertTrue(t, ipIs4(ip4), "ip is not ipv4")
    assertFalse(t, ipIs4(ip6), "ip is ipv4")
}

func Test_ipIs6(t *testing.T) {
    var ip4 = net.ParseIP("1.1.1.1")
    var ip6 = net.ParseIP("2001:4860:0:2001::68")

    assertTrue(t, ipIs6(ip6), "ip is ipv6")
    assertFalse(t, ipIs6(ip4), "ip is not ipv6")
}


func Test_maskIs4(t *testing.T) {
    var mask4 = net.CIDRMask(8, 32)
    var mask6 = net.CIDRMask(8, 128)

    assertTrue(t, maskIs4(mask4), "mask is ipv4")
    assertFalse(t, maskIs4(mask6), "mask is not ipv4")
}

func Test_maskIs6(t *testing.T) {
    var mask4 = net.CIDRMask(8, 32)
    var mask6 = net.CIDRMask(8, 128)

    assertTrue(t, maskIs6(mask6), "mask is ipv6")
    assertFalse(t, maskIs6(mask4), "mask is not ipv6")
}


func Test_ipmaskAsString4(t *testing.T) {
    var mask = net.CIDRMask(24, 32)

    assertStrEq(t, "255.255.255.0", ipmaskAsString4(mask), "failed to format mask")
}


func Test_ipIsLesser4(t *testing.T) {
    var a = net.ParseIP("10.0.2.15")
    var d = net.ParseIP("10.0.2.16")

    var g = net.ParseIP("10.0.3.14")
    var h = net.ParseIP("10.0.3.15")

    var m = net.ParseIP("10.1.2.14")
    var n = net.ParseIP("10.1.2.15")

    var v = net.ParseIP("11.0.2.14")
    var w = net.ParseIP("11.0.2.15")


    assertFalse(t, ipIsLesser(a, a), "a == a")
    assertFalse(t, ipIsLesser(d, d), "d == d")

    assertTrue(t, ipIsLesser(a, d), "a < d")
    assertFalse(t, ipIsLesser(d, a), "d !< a")

    assertTrue(t, ipIsLesser(a, g), "a < g")
    assertTrue(t, ipIsLesser(a, h), "a < h")
    assertFalse(t, ipIsLesser(g, a), "g !< a")
    assertFalse(t, ipIsLesser(h, a), "h !< a")

    assertTrue(t, ipIsLesser(a, m), "a < m")
    assertTrue(t, ipIsLesser(a, n), "a < n")
    assertFalse(t, ipIsLesser(m, a), "m !< a")
    assertFalse(t, ipIsLesser(n, a), "n !< a")

    assertTrue(t, ipIsLesser(a, v), "a < v")
    assertTrue(t, ipIsLesser(a, w), "a < w")
    assertFalse(t, ipIsLesser(v, a), "v !< a")
    assertFalse(t, ipIsLesser(w, a), "w !< a")
}

func Test_ipIsLesser6(t *testing.T) {
    var a = net.ParseIP("2001:4860:0:2001::68")
    var b = net.ParseIP("2001:4860:0:2001::69")
    var c = net.ParseIP("2001:4860:0:2002::68")
    var d = net.ParseIP("2001:4860:0:2002::67")

    assertFalse(t, ipIsLesser(a, a), "a == a")

    assertTrue(t, ipIsLesser(a, b), "a < b")
    assertFalse(t, ipIsLesser(b, a), "b !< a")

    assertTrue(t, ipIsLesser(a, c), "a < c")
    assertFalse(t, ipIsLesser(c, a), "c !< a")

    assertTrue(t, ipIsLesser(a, d), "a < d")
    assertFalse(t, ipIsLesser(d, a), "d !< a")
}


func Test_applyMask4(t *testing.T) {
    var ip4 = net.ParseIP("122.204.201.241")
    var mask4zero = net.IPv4Mask(255, 255, 255, 255)
    var mask4one = net.IPv4Mask(255, 255, 255, 0)
    var mask4two = net.IPv4Mask(255, 255, 0, 0)
    var mask4onehalf = net.IPv4Mask(255, 128, 0, 0)

    assertStrEq(t, applyMask(&ip4, &mask4zero).IP.String(),
                    "122.204.201.241", "ip4 masking failed")
    assertStrEq(t, applyMask(&ip4, &mask4one).IP.String(),
                    "122.204.201.0", "ip4 masking failed")
    assertStrEq(t, applyMask(&ip4, &mask4two).IP.String(),
                    "122.204.0.0", "ip4 masking failed")
    assertStrEq(t, applyMask(&ip4, &mask4onehalf).IP.String(),
                    "122.128.0.0", "ip4 masking failed")
}

func Test_applyMask6(t *testing.T) {
    var ip6 = net.ParseIP("2001:4860:0:2002::fe:248")
    var mask6zero = net.CIDRMask(128, 128)
    var mask6one = net.CIDRMask(120, 128)
    var mask6two = net.CIDRMask(112, 128)
    var mask6three = net.CIDRMask(106, 128)
    var mask6four = net.CIDRMask(98, 128)

    var mask6twentynine = net.CIDRMask(24, 128)
    var mask6thirty = net.CIDRMask(16, 128)
    var mask6thirtyone = net.CIDRMask(8, 128)

    assertStrEq(t, applyMask(&ip6, &mask6zero).IP.String(),
                    "2001:4860:0:2002::fe:248", "ip6 masking failed")
    assertStrEq(t, applyMask(&ip6, &mask6one).IP.String(),
                    "2001:4860:0:2002::fe:200", "ip6 masking failed")
    assertStrEq(t, applyMask(&ip6, &mask6two).IP.String(),
                    "2001:4860:0:2002::fe:0", "ip6 masking failed")
    assertStrEq(t, applyMask(&ip6, &mask6three).IP.String(),
                    "2001:4860:0:2002::c0:0", "ip6 masking failed")
    assertStrEq(t, applyMask(&ip6, &mask6four).IP.String(),
                    "2001:4860:0:2002::", "ip6 masking failed")

    assertStrEq(t, applyMask(&ip6, &mask6twentynine).IP.String(),
                    "2001:4800::", "ip6 masking failed")
    assertStrEq(t, applyMask(&ip6, &mask6thirty).IP.String(),
                    "2001::", "ip6 masking failed")
    assertStrEq(t, applyMask(&ip6, &mask6thirtyone).IP.String(),
                    "2000::", "ip6 masking failed")
}


func Test_maskAsIpToIPMask(t *testing.T) {
    var ip4half = net.ParseIP("255.255.255.127")
    var ip4one = net.ParseIP("255.255.255.0")
    var ip6quarter = net.ParseIP("ffff:ffff::")
    var ip6half = net.ParseIP("ffff:ffff:ffff:ffff::")

    assertStrEq(t, "ffffff7f", maskAsIpToIPMask(&ip4half).String(), "wrong ipv4 mask")
    assertStrEq(t, "ffffff00", maskAsIpToIPMask(&ip4one).String(), "wrong ipv4 mask")

    assertStrEq(t, "ffffffff000000000000000000000000",
                    maskAsIpToIPMask(&ip6quarter).String(), "wrong ipv6 mask")
    assertStrEq(t, "ffffffffffffffff0000000000000000",
                    maskAsIpToIPMask(&ip6half).String(), "wrong ipv6 mask")
}


func Test_ipnetMaskAsIP(t *testing.T) {
    var _, ipnet4, _ = net.ParseCIDR("144.124.153.123/24")
    var _, ipnet6fst, _ = net.ParseCIDR("2001:4860:0:2002::/64")
    var _, ipnet6snd, _ = net.ParseCIDR("2001:4860:0:2002::/96")

    assertStrEq(t, "255.255.255.0", ipnetMaskAsIP(ipnet4).String(),
                                    "ipv4 mask extraction failed")
    assertStrEq(t, "ffff:ffff:ffff:ffff::",
                    ipnetMaskAsIP(ipnet6fst).String(), "ipv6 mask extraction failed")
    assertStrEq(t, "ffff:ffff:ffff:ffff:ffff:ffff::",
                    ipnetMaskAsIP(ipnet6snd).String(), "ipv6 mask extraction failed")
}
