package gosha

func (e *Exec) Run() error {
	for _, cmd := range e.Commands {
		if err := runCmd(e.Scanner, cmd); err != nil {
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
		go syncRunCmd(e.Scanner, cmd, errCh)
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
