package xproc

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/michurin/cnbot/ctxlog"
	"github.com/michurin/cnbot/xlog"
)

type Cmd struct {
	InterruptDelay time.Duration
	KillDelay      time.Duration
	Env            []string
	Command        string
	Cwd            string
}

func killGrp(ctx context.Context, pid int, sig syscall.Signal) {
	// We do not consider error as critical because the process could
	// disappear by its own. It is not easy to identify error in this case.
	// For example you can get ESRCH (0x3) that doesn't support by syscall.Errno.Is().
	pgid, err := syscall.Getpgid(pid) // not cmd.SysProcAttr.Pgid
	if err != nil {
		xlog.L(ctx, fmt.Errorf("kill: getpgid: %w", err))
		return
	}
	err = syscall.Kill(-pgid, sig) // minus
	if err != nil {
		xlog.L(ctx, fmt.Errorf("kill: kill %d: %w", -pgid, err))
		return
	}
}

// Note: don't use ctx for timeouts
// Note: be careful with env argument, side effects are possible
func (c *Cmd) Run(
	ctx context.Context,
	args []string,
	env []string,
) (
	[]byte,
	error,
) {
	// setup cmd
	cmd := exec.Command(c.Command, args...) // we don't use CommandContext here because it kills only process, not group
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Dir = c.Cwd
	cmd.Env = append(env, c.Env...)
	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	var errBuffer bytes.Buffer
	cmd.Stderr = &errBuffer

	xlog.L(ctx, fmt.Sprintf("starting %s %v", c.Command, args)) // TODO put command to context?

	err := cmd.Start() // start command synchronously
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "start: %w", err)
	}
	ctx = ctxlog.Add(ctx, "pid", cmd.Process.Pid)

	done := make(chan struct{})
	intBound := time.NewTimer(c.InterruptDelay)
	killBound := time.NewTimer(c.InterruptDelay + c.KillDelay)
	defer func() {
		intBound.Stop()
		killBound.Stop()
		close(done)
	}()
	go func() {
		for {
			select {
			case <-done: // it has to appear before kill sections to catch stat errors
				return
			case <-ctx.Done(): // urgent exit, we doesn't even wait for process finalization
				xlog.L(ctx, "Exec terminated by context")
				killGrp(ctx, cmd.Process.Pid, syscall.SIGKILL)
				return
			case <-intBound.C:
				killGrp(ctx, cmd.Process.Pid, syscall.SIGINT) // Not all OS support SIGTERM
			case <-killBound.C:
				killGrp(ctx, cmd.Process.Pid, syscall.SIGKILL)
			}
		}
	}()
	err = cmd.Wait()
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "wait: %w", err)
	}

	errMsg := []string(nil)
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		errMsg = append(errMsg, fmt.Sprintf("exit code: %d", exitCode))
	}
	errStr := errBuffer.String()
	if errStr != "" {
		xlog.L(ctx, fmt.Errorf("stderr: %s", errStr)) // TODO consider as error?
		// errMsg = append(errMsg, fmt.Sprintf("stderr: %q", errStr))
	}
	outBytes := outBuffer.Bytes()
	if errMsg == nil {
		return outBytes, nil
	}
	errMsg = append(errMsg, fmt.Sprintf("stdout: %q", string(outBytes))) // TODO trim?
	return nil, ctxlog.Errorfx(ctx, strings.Join(errMsg, "; "))
}
