package main

import "testing"


func assertIntEq(t *testing.T, a, b int, msg string) {
    if a != b {
        t.Errorf("Integers do not match: %d != %d [%s]", a, b, msg)
    }
}

func assertPtrIsNull(t *testing.T, a interface{}, msg string) {
    if a != nil {
        t.Errorf("Pointer is not nil: %d, {%q} [%s]", a, a, msg)
    }
}

func assertStrEq(t *testing.T, a, b string, msg string) {
    if a != b {
        t.Errorf("Strings do not match: '%s' != '%s' [%s]", a, b, msg)
    }
}
