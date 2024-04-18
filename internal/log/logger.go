package log

import (
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/pkg/logging"
	"sync"
)

var (
	loggerInstance *logging.Logger
	loggerOnce     sync.Once
)

func GetLogger() *logging.Logger {
	loggerOnce.Do(func() {
		mode := logging.ProdMode

		switch cfg := config.GetConfig(); cfg.Project.Mode {
		case "local":
			mode = logging.LocalMode
		case "dev":
			mode = logging.DevMode
		}

		loggerInstance = logging.GetLogger(mode)
	})

	return loggerInstance
}
