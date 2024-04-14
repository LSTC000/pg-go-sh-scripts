package usecase

import (
	"context"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/internal/service"

	uuid "github.com/satori/go.uuid"
)

type (
	IBashLogUseCase interface {
		GetBashLogPaginationPageByBashId(bashId uuid.UUID, paginationParams schema.PaginationParams) (schema.PaginationPage[*model.BashLog], error)
	}

	BashLogUseCase struct {
		service    service.IBashLogService
		httpErrors *config.HTTPErrors
	}
)

func (u *BashLogUseCase) GetBashLogPaginationPageByBashId(bashId uuid.UUID, paginationParams schema.PaginationParams) (schema.PaginationPage[*model.BashLog], error) {
	var bashLogPaginationPage schema.PaginationPage[*model.BashLog]

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
