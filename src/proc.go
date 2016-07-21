package main

import "bytes"
import "fmt"
import "os"
import "os/exec"
import "syscall"


type ProcessManager struct {
    exe string
    args []string
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
    var cmd = exec.Command(mgr.exe, mgr.args...)

    // Set up buffers for stdout, stderr
    var errBuffer bytes.Buffer
    var outBuffer bytes.Buffer
    cmd.Stderr = &errBuffer
    cmd.Stdout = &outBuffer

    // Invoke the cmd
    os.Setenv("LC_ALL", "C")  // Make sure we don't get localized output
    var err = cmd.Run()
    if err != nil {
        return ProcessResult{
            err: fmt.Errorf("Failed to run %s %s: %q", mgr.exe, mgr.args, err),
        }
    }

    // Capture stdout, stderr
    var stderr = errBuffer.String()
    var stdout = outBuffer.String()
    var waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)

    // Construct a result
    var res = ProcessResult{
        exitCode: waitStatus.ExitStatus(),
        stderr: stderr,
        stdout: stdout,
        err: err,
    }

    return res
}
