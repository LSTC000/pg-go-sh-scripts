package server

import (
	"errors"
	"fmt"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/log"
	"pg-sh-scripts/internal/schema"

	"github.com/gin-gonic/gin"
)

func getLogMiddleware() gin.HandlerFunc {
	return gin.Logger()
}

func getRecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var httpErr *schema.HTTPError

				logger := log.GetLogger()
				logger.Error(fmt.Sprintf("Unknown error: %v", err))

				httpErrors := config.GetHTTPErrors()
				errors.As(httpErrors.Internal, &httpErr)

				ctx.JSON(httpErr.HTTPCode, httpErr)
			}
		}()
		ctx.Next()
	}
}

func setServeMiddleware(r *gin.Engine) {
	r.Use(getLogMiddleware(), getRecoveryMiddleware())
}
