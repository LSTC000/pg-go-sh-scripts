package gosha

const (
	execOperator = "/bin/bash"
)

type (
	IExec interface {
		Run() error
		SyncRun() []error
	}

	Exec struct {
		Scanner  IScanner
		Commands []ICmd
	}
)

func (e *Exec) Run() error {
	for _, cmd := range e.Commands {
		if err := cmd.run(e.Scanner); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exec) SyncRun() []error {
	commandsCount := len(e.Commands)

	errPool := make([]error, 0, commandsCount)
	errCh := make(chan error)
	defer close(errCh)

	for _, cmd := range e.Commands {
		go cmd.syncRun(e.Scanner, errCh)
	}

	for i := 0; i < commandsCount; i++ {
		if err := <-errCh; err != nil {
			errPool = append(errPool, err)
		}
	}

	if len(errPool) > 0 {
		return errPool
	}

	return nil
}

func GetExec(scanner IScanner, commands []ICmd) IExec {
	return &Exec{
		Scanner:  scanner,
		Commands: commands,
	}
}
