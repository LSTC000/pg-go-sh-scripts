package gosha

const (
	execOperator = "/bin/bash"
)

type (
	IExec interface {
		Run(IScanner, []ICmd) error
		SyncRun(IScanner, []ICmd) []error
	}

	Exec struct{}
)

func (e *Exec) Run(scanner IScanner, commands []ICmd) error {
	for _, cmd := range commands {
		if err := cmd.run(scanner); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exec) SyncRun(scanner IScanner, commands []ICmd) []error {
	commandsCount := len(commands)

	errPool := make([]error, 0, commandsCount)
	errCh := make(chan error)
	defer close(errCh)

	for _, cmd := range commands {
		go cmd.syncRun(scanner, errCh)
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

func GetExec() IExec {
	return &Exec{}
}
