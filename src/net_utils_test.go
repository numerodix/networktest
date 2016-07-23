package main

import "net"
import "testing"


func Test_ipIsLesser(t *testing.T) {
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
