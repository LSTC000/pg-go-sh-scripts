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
		Commands []Cmd
	}
)

func GetExec(scanner IScanner, commands []Cmd) IExec {
	return &Exec{
		Scanner:  scanner,
		Commands: commands,
	}
}
