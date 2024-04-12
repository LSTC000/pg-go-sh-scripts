package server

import (
	"fmt"
	"net/http"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/log"

	"github.com/gin-gonic/gin"
)

func getLogMiddleware() gin.HandlerFunc {
	return gin.Logger()
}

func getRecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger := log.GetLogger()
				logger.Error(fmt.Sprintf("Unknown Error: %v", err))
				httpErrors := config.GetHTTPErrors()
				ctx.JSON(http.StatusInternalServerError, httpErrors.Internal)
			}
		}()
		ctx.Next()
	}
}

func setServeMiddleware(r *gin.Engine) {
	r.Use(getLogMiddleware(), getRecoveryMiddleware())
}
