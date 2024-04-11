package common

import (
	"pg-sh-scripts/pkg/gosha"
)

func GetGoshaExec(scanner gosha.IScanner, commands []gosha.Cmd) gosha.IExec {
	return gosha.GetExec(scanner, commands)
}
