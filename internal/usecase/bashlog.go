package usecase

import (
	"context"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/service"

	uuid "github.com/satori/go.uuid"
)

type (
	IBashLogUseCase interface {
		GetBashLogListByBashId(bashId uuid.UUID) ([]*model.BashLog, error)
	}

	BashLogUseCase struct {
		service    service.IBashLogService
		httpErrors *config.HTTPErrors
	}
)

func (u *BashLogUseCase) GetBashLogListByBashId(bashId uuid.UUID) ([]*model.BashLog, error) {
	bashService := service.GetBashService()
	_, err := bashService.GetOneById(context.Background(), bashId)
	if err != nil {
		return nil, u.httpErrors.BashGet
	}

	bashLogList, err := u.service.GetAllByBashId(context.Background(), bashId)
	if err != nil {
		return nil, u.httpErrors.BashLogGetListByBashId
	}

	return bashLogList, nil
}

func GetBashLogUseCase() IBashLogUseCase {
	return &BashLogUseCase{
		service:    service.GetBashLogService(),
		httpErrors: config.GetHTTPErrors(),
	}
}
