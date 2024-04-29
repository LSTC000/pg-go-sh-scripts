package usecase

import (
	"context"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/service"
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/pkg/sql/pagination"

	uuid "github.com/satori/go.uuid"
)

//go:generate mockgen -source=./bashlog.go  -destination=./mock/bashlog.go

type (
	IBashLogUseCase interface {
		GetBashLogPaginationPageByBashId(
			bashId uuid.UUID,
			paginationParams pagination.LimitOffsetParams,
		) (alias.BashLogLimitOffsetPage, error)
	}

	BashLogUseCase struct {
		service     service.IBashLogService
		bashService service.IBashService
		httpErrors  *config.HTTPErrors
	}
)

func (u *BashLogUseCase) GetBashLogPaginationPageByBashId(
	bashId uuid.UUID,
	paginationParams pagination.LimitOffsetParams,
) (alias.BashLogLimitOffsetPage, error) {
	var bashLogPaginationPage alias.BashLogLimitOffsetPage

	_, err := u.bashService.GetOneById(context.Background(), bashId)
	if err != nil {
		return bashLogPaginationPage, u.httpErrors.BashDoesNotExists
	}

	bashLogPaginationPage, err = u.service.GetPaginationPageByBashId(
		context.Background(),
		bashId,
		paginationParams,
	)
	if err != nil {
		return bashLogPaginationPage, u.httpErrors.BashLogGetPaginationPageByBashId
	}

	return bashLogPaginationPage, nil
}

func GetBashLogUseCase() IBashLogUseCase {
	return &BashLogUseCase{
		service:     service.GetBashLogService(),
		bashService: service.GetBashService(),
		httpErrors:  config.GetHTTPErrors(),
	}
}
