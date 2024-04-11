package server

import (
	bashV1 "pg-sh-scripts/internal/api/v1/bash"
	bashLogV1 "pg-sh-scripts/internal/api/v1/bashlog"
	"pg-sh-scripts/internal/config"

	"github.com/gin-gonic/gin"
)

func setV1Handlers(r *gin.Engine, cfg *config.Config) {
	rg := r.Group(cfg.Api.Prefix)

	bashV1Handler := bashV1.GetHandler()
	bashV1Handler.Register(rg)

	bashLogV1Handler := bashLogV1.GetHandler()
	bashLogV1Handler.Register(rg)
}
