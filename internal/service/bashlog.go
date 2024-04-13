package service

import (
	"context"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/repo"

	uuid "github.com/satori/go.uuid"
)

type (
	IBashLogService interface {
		GetAllByBashId(ctx context.Context, bashId uuid.UUID) ([]*model.BashLog, error)
		Create(ctx context.Context, dto dto.CreateBashLogDTO) (*model.BashLog, error)
	}

	BashLogService struct {
		repository repo.IBashLogRepository
	}
)

func (s *BashLogService) GetAllByBashId(ctx context.Context, bashId uuid.UUID) ([]*model.BashLog, error) {
	bashLogList, err := s.repository.GetAllByBashId(ctx, bashId)
	if err != nil {
		return nil, err
	}
	return bashLogList, nil
}

func (s *BashLogService) Create(ctx context.Context, dto dto.CreateBashLogDTO) (*model.BashLog, error) {
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
