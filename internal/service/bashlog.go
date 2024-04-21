package service

import (
	"context"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/repo"
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/pkg/sql/pagination"

	uuid "github.com/satori/go.uuid"
)

type (
	IBashLogService interface {
		GetPaginationPageByBashId(ctx context.Context, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) (alias.BashLogLimitOffsetPage, error)
		Create(ctx context.Context, dto dto.CreateBashLog) (*model.BashLog, error)
	}

	BashLogService struct {
		repository repo.IBashLogRepository
	}
)

func (s *BashLogService) GetPaginationPageByBashId(ctx context.Context, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) (alias.BashLogLimitOffsetPage, error) {
	bashLogPaginationPage, err := s.repository.GetPaginationPageByBashId(ctx, bashId, paginationParams)
	if err != nil {
		return bashLogPaginationPage, err
	}
	return bashLogPaginationPage, nil
}

func (s *BashLogService) Create(ctx context.Context, dto dto.CreateBashLog) (*model.BashLog, error) {
	bashLog, err := s.repository.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	return bashLog, nil
}

func GetBashLogService() IBashLogService {
	return &BashLogService{
		repository: repo.GetPgBashLogRepository(),
	}
}
