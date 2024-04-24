package api

import (
	"errors"
	"fmt"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/log"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/pkg/logging"
)

//go:generate mockgen -source=./helper.go  -destination=./mock/helper.go

type (
	IHelper interface {
		ParseError(err error) *schema.HTTPError
	}

	Helper struct {
		logger     *logging.Logger
		httpErrors *config.HTTPErrors
	}
)

func (e *Helper) ParseError(err error) *schema.HTTPError {
	var httpErr *schema.HTTPError

	if errors.As(err, &httpErr) {
		e.logger.Error(fmt.Sprintf("Service error: %v", err))
	} else {
		e.logger.Error(fmt.Sprintf("Unknown error: %v", err))
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
