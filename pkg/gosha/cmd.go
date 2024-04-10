package gosha

import (
	"context"
	"os/exec"
	"time"
)

type Cmd struct {
	Title   string
	Path    string
	Timeout time.Duration
}

func runCmd(scanner IScanner, cmd Cmd) error {
	var cmdExec *exec.Cmd

	cmdPath := cmd.Path
	cmdTimeout := cmd.Timeout

	if cmdTimeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
		defer cancel()
		cmdExec = exec.CommandContext(ctx, execOperator, cmdPath)
	} else {
		cmdExec = exec.Command(execOperator, cmdPath)
	}

	stdout, err := cmdExec.StdoutPipe()
	if err != nil {
		return GetExecErr(cmdPath, ErrFmt(stdoutErrGroup, err))
	}

	if err = cmdExec.Start(); err != nil {
		return GetExecErr(cmdPath, ErrFmt(startExecErrGroup, err))

	}

	if err = scanner.Scan(stdout, cmd); err != nil {
		return GetExecErr(cmdPath, ErrFmt(scanErrGroup, err))
	}

	if err = cmdExec.Wait(); err != nil {
		return GetExecErr(cmdPath, ErrFmt(waitExecErrGroup, err))
	}

	return nil
}

func syncRunCmd(scanner IScanner, cmd Cmd, ch chan<- error) {
	ch <- runCmd(scanner, cmd)
}
