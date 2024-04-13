package api

import (
	"errors"
	"fmt"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/log"
	"pg-sh-scripts/internal/schema"

	"github.com/gin-gonic/gin"
)

func RaiseError(ctx *gin.Context, err error) {
	var httpErr *schema.HTTPError

	logger := log.GetLogger()

	if errors.As(err, &httpErr) {
		logger.Error(fmt.Sprintf("Path: %s Error: %v", ctx.FullPath(), err))
		ctx.JSON(httpErr.HTTPCode, httpErr)
	} else {
		httpErrors := config.GetHTTPErrors()
		logger.Error(fmt.Sprintf("Path: %s Unknown Error: %v", ctx.FullPath(), err))
		errors.As(httpErrors.Internal, &httpErr)
		ctx.JSON(httpErr.HTTPCode, httpErr)
	}
}
