package server

import (
	v1 "pg-sh-scripts/internal/api/v1"
	"pg-sh-scripts/internal/config"

	"github.com/gin-gonic/gin"
)

func setV1Handlers(r *gin.Engine, cfg *config.Config) {
	rg := r.Group(cfg.Api.Prefix)

	bashV1Handler := v1.GetBashHandler()
	bashV1Handler.Register(rg)

	bashLogV1Handler := v1.GetBashLogHandler()
	bashLogV1Handler.Register(rg)
}
