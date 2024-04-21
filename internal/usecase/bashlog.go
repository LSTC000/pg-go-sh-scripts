package usecase

import (
	"context"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/service"
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/pkg/sql/pagination"

	uuid "github.com/satori/go.uuid"
)

type (
	IBashLogUseCase interface {
		GetBashLogPaginationPageByBashId(bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) (alias.BashLogLimitOffsetPage, error)
	}

	BashLogUseCase struct {
		service    service.IBashLogService
		httpErrors *config.HTTPErrors
	}
)

func (u *BashLogUseCase) GetBashLogPaginationPageByBashId(bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) (alias.BashLogLimitOffsetPage, error) {
	var bashLogPaginationPage pagination.LimitOffsetPage[*model.BashLog]

	bashService := service.GetBashService()
	_, err := bashService.GetOneById(context.Background(), bashId)
	if err != nil {
		return bashLogPaginationPage, u.httpErrors.BashDoesNotExists
	}

	bashLogPaginationPage, err = u.service.GetPaginationPageByBashId(context.Background(), bashId, paginationParams)
	if err != nil {
		return bashLogPaginationPage, u.httpErrors.BashLogGetPaginationPageByBashId
	}

	return bashLogPaginationPage, nil
}

func GetBashLogUseCase() IBashLogUseCase {
	return &BashLogUseCase{
		service:    service.GetBashLogService(),
		httpErrors: config.GetHTTPErrors(),
	}
}
