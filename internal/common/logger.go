package common

import (
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/pkg/logging"
	"sync"
)

func GetLogger() *logging.Logger {
	var once sync.Once

	mode := logging.ProdMode

	once.Do(func() {
		switch cfg := config.GetConfig(); cfg.Project.Mode {
		case "local":
			mode = logging.LocalMode
		case "dev":
			mode = logging.DevMode
		}
	})

	return logging.GetLogger(mode)
}
