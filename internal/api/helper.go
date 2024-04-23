package api

import (
	"errors"
	"fmt"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/log"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/pkg/logging"

	"github.com/gin-gonic/gin"
)

type (
	IHelper interface {
		ParseError(*gin.Context, error) *schema.HTTPError
	}

	Helper struct {
		logger     *logging.Logger
		httpErrors *config.HTTPErrors
	}
)

func (e *Helper) ParseError(c *gin.Context, err error) *schema.HTTPError {
	var httpErr *schema.HTTPError

	if errors.As(err, &httpErr) {
		e.logger.Error(fmt.Sprintf("Path: %s Error: %v", c.FullPath(), err))
	} else {
		e.logger.Error(fmt.Sprintf("Path: %s Unknown Error: %v", c.FullPath(), err))
		errors.As(e.httpErrors.Internal, &httpErr)
	}

	return httpErr
}

func GetHelper() IHelper {
	return &Helper{
		logger:     log.GetLogger(),
		httpErrors: config.GetHTTPErrors(),
	}
}
