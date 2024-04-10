package common

import (
	"pg-sh-scripts/pkg/gosha"
)

func GetGoshaExec(commands []gosha.Cmd) gosha.IExec {
	return gosha.GetExec(gosha.GetDefaultScanner(), commands)
}
