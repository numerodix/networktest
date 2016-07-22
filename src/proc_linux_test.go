package main

//import "fmt"
import "strings"
import "testing"


func Test_ProcMgr_false(t *testing.T) {
    var mgr = ProcMgr("false")
    var res = mgr.run()

    assertStrEq(t, "", res.stderr, "stderr does not match")
    assertStrEq(t, "", res.stdout, "stdout does not match")

    assertPtrNotNull(t, res.err, "err does not match")
    assertStrEq(t, "exit status 1", res.err.Error(), "err does not match")
}

func Test_ProcMgr_ls(t *testing.T) {
    var mgr = ProcMgr("ls", "/dyhfi8345rh")
    var res = mgr.run()

    // See that we don't get localized output
    assertStrEq(t, "ls: cannot access /dyhfi8345rh: No such file or directory",
                    strings.TrimSpace(res.stderr), "stderr does not match")
    assertStrEq(t, "", strings.TrimSpace(res.stdout), "stdout does not match")

    assertPtrNotNull(t, res.err, "err does not match")
    assertStrEq(t, "exit status 2", res.err.Error(), "err does not match")
}

func Test_ProcMgr_true(t *testing.T) {
    var mgr = ProcMgr("true")
    var res = mgr.run()

    assertStrEq(t, "", res.stderr, "stderr does not match")
    assertStrEq(t, "", res.stdout, "stdout does not match")

    assertPtrIsNull(t, res.err, "err does not match")
}

func Test_ProcMgr_uname(t *testing.T) {
    var mgr = ProcMgr("uname")
    var res = mgr.run()

    assertStrEq(t, "", res.stderr, "stderr does not match")
    assertStrEq(t, "Linux", strings.TrimSpace(res.stdout), "stdout does not match")

    assertPtrIsNull(t, res.err, "err does not match")
}
