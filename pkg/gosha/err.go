package gosha

import "fmt"

type (
	ErrGroup string

	ExecErr struct {
		Title  string
		Path   string
		Detail string
	}
)

const (
	stdoutErrGroup    ErrGroup = "stdout"
	startExecErrGroup ErrGroup = "start execute"
	waitExecErrGroup  ErrGroup = "wait execute"
	scanErrGroup      ErrGroup = "scan"
)

func (e *ExecErr) Error() string {
	return fmt.Sprintf("gosha was shocked - %s", e.Detail)
}

func ErrFmt(group ErrGroup, err error) string {
	return fmt.Sprintf("[%s] error: %s", group, err)
}

func GetExecErr(cmd Cmd, detail string) error {
	return &ExecErr{
		Title:  cmd.Title,
		Path:   cmd.Path,
		Detail: detail,
	}
}
