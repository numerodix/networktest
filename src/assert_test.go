package main

import "fmt"
import "runtime"
import "strings"
import "testing"


// Peek up the call stack to find the name of the calling function
func _getFuncPtrVal() uintptr {
    // frame 0 : assertXXX
    // frame 1 : TestXXX <- the one we want
    // frame 2 : testing.tRunner
    var frameCnt = 3

    // Extract the top few frames
    var pc = make([]uintptr, frameCnt)
    runtime.Callers(frameCnt, pc)

    return pc[1]
}

// Return the function name for the calling function
func _getFuncName() string {
    var ptrVal = _getFuncPtrVal()

    // Look up the function object
    var fun = runtime.FuncForPC(ptrVal)

    // returns: _/home/user/src/gohavenet/src.TestProcMgr_False
    var funcNamePath = fun.Name()

    // Split on the slash and return just the func name
    var pathElems = strings.Split(funcNamePath, "/")
    var index = len(pathElems) - 1
    if index < 0 {
        index = 0
    }

    return pathElems[index]
}

// Return /path/to/file:12 for the calling function
func _getFuncFileLine() string {
    var ptrVal = _getFuncPtrVal()

    // Look up the function object
    var fun = runtime.FuncForPC(ptrVal)
    var file, line = fun.FileLine(ptrVal)

    var fileLine = fmt.Sprintf("%s:%d", file, line)

    return fileLine
}



func assertFalse(t *testing.T, b bool, msg string) {
    if b {
        t.Errorf("%s: Was not false: %s != false [%s] at %s",
                 _getFuncName(), b, msg, _getFuncFileLine())
    }
}

func assertTrue(t *testing.T, b bool, msg string) {
    if !b {
        t.Errorf("%s: Was not true: %s != true [%s] at %s",
                 _getFuncName(), b, msg, _getFuncFileLine())
    }
}

func assertIntEq(t *testing.T, a, b int, msg string) {
    if a != b {
        t.Errorf("%s: Integers do not match: %d != %d [%s] at %s",
                 _getFuncName(), a, b, msg, _getFuncFileLine())
    }
}

func assertPtrIsNull(t *testing.T, a interface{}, msg string) {
    if a != nil {
        t.Errorf("%s: Pointer is not nil: %d, {%q} [%s] at %s",
                 _getFuncName(), a, a, msg, _getFuncFileLine())
    }
}

func assertPtrNotNull(t *testing.T, a interface{}, msg string) {
    if a == nil {
        t.Errorf("%s: Pointer is nil: %d, [%s] at %s",
                 _getFuncName(), a, msg, _getFuncFileLine())
    }
}

func assertStrEq(t *testing.T, a, b string, msg string) {
    if a != b {
        t.Errorf("%s: Strings do not match: '%s' != '%s' [%s] at %s",
                 _getFuncName(), a, b, msg, _getFuncFileLine())
    }
}

func assertPtrEq(t *testing.T, a, b interface{}, msg string) {
    if a != b {
        t.Errorf("%s: Interfaces do not match: %s != %s [%s] at %s",
                 _getFuncName(), a, b, msg, _getFuncFileLine())
    }
}
