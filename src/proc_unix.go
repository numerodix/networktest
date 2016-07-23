package main

import "bytes"
//import "fmt"
import "os"
import "os/exec"
import "time"


type ProcessManager struct {
    exe string
    args []string
    timeoutMs int
}


func ProcMgr(exe string, args... string) ProcessManager {
    var mgr = ProcessManager{
        exe: exe,
        args: args,
    }
    return mgr
}


func (mgr *ProcessManager) run() ProcessResult {
    var res ProcessResult

    if mgr.timeoutMs > 0 {
        res = mgr.runWithTimeout(mgr.timeoutMs)
    } else {
        res = mgr.runStandard()
    }

    return res
}


func (mgr *ProcessManager) runStandard() ProcessResult {
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

    // Capture stdout, stderr
    var stderr = errBuffer.String()
    var stdout = outBuffer.String()

    // Construct a result
    var res = ProcessResult{
        stderr: stderr,
        stdout: stdout,
        err: err,
    }

    return res
}


func (mgr *ProcessManager) runWithTimeout(timeoutMs int) ProcessResult {
    // Construct the cmd
    var cmd = exec.Command(mgr.exe, mgr.args...)

    // Set up buffers for stdout, stderr
    var errBuffer bytes.Buffer
    var outBuffer bytes.Buffer
    cmd.Stderr = &errBuffer
    cmd.Stdout = &outBuffer

    // Invoke the cmd
    os.Setenv("LC_ALL", "C")  // Make sure we don't get localized output
    var err = cmd.Start()

    var elapsedMs = 0
    for {
        // We reached the timeout, so kill the process and return an error
        if timeoutMs > 0 && elapsedMs >= timeoutMs {
            err = cmd.Process.Kill()
            break
        }

        // Use the pid to detect if the process exists
        var _, err = os.FindProcess(cmd.Process.Pid)

        // If it doesn't:
        // - it's still forking (not yet running)
        // - or it exited already (we check stdout to see if it's non-empty)
        if err == nil && (outBuffer.String() != "" || errBuffer.String() != "") {
            break
        }

        // The process hasn't finished yet, so wait a while and loop around
        time.Sleep(50 * time.Millisecond)
        elapsedMs += 50
    }

    // Capture stdout, stderr
    var stderr = errBuffer.String()
    var stdout = outBuffer.String()

    // Construct a result
    var res = ProcessResult{
        stderr: stderr,
        stdout: stdout,
        err: err,
    }

    return res
}
