package gosha

import (
	"context"
	"os/exec"
	"time"
)

type (
	ICmd interface {
		run(IScanner) error
		syncRun(IScanner, chan<- error)
	}

	Cmd struct {
		Title   string
		Path    string
		Timeout time.Duration
	}
)

func (c *Cmd) run(scanner IScanner) error {
	var cmdExec *exec.Cmd

	cmdPath := c.Path
	cmdTimeout := c.Timeout

	if cmdTimeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
		defer cancel()
		cmdExec = exec.CommandContext(ctx, execOperator, cmdPath)
	} else {
		cmdExec = exec.Command(execOperator, cmdPath)
	}

	stdout, err := cmdExec.StdoutPipe()
	if err != nil {
		return GetExecErr(c, ErrFmt(stdoutErrGroup, err))
	}

	if err = cmdExec.Start(); err != nil {
		return GetExecErr(c, ErrFmt(startExecErrGroup, err))

	}

	if err = scanner.Scan(stdout, c); err != nil {
		return GetExecErr(c, ErrFmt(scanErrGroup, err))
	}

	if err = cmdExec.Wait(); err != nil {
		return GetExecErr(c, ErrFmt(waitExecErrGroup, err))
	}

	return nil
}

func (c *Cmd) syncRun(scanner IScanner, ch chan<- error) {
	ch <- c.run(scanner)
}
