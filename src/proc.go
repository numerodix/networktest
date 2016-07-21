package main

import "bytes"
import "fmt"
import "os"
import "os/exec"


type ProcessManager struct {
    exe string
    args []string
}

type ProcessResult struct {
    exitCode int
    stdout string
    stderr string
    err error
}

func ProcMgr(exe string, args... string) ProcessManager {
    var mgr = ProcessManager{
        exe: exe,
        args: args,
    }
    return mgr
}

func (mgr *ProcessManager) run() ProcessResult {
    // Construct the cmd
    cmd := exec.Command(mgr.exe, mgr.args...)
    var outBuf bytes.Buffer
    cmd.Stdout = &outBuf

    // Invoke the cmd
    os.Setenv("LC_ALL", "C")
    err := cmd.Run()
    if err != nil {
        return ProcessResult{
            err: fmt.Errorf("Failed to run %s %s: %q", mgr.exe, mgr.args, err),
        }
    }

    // Parse the output into lines
    var stdout = outBuf.String()

    var res = ProcessResult{
        stdout: stdout,
        err: err,
    }

    return res
}
