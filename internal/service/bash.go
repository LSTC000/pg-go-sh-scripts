package service

import (
	"context"
	"pg-sh-scripts/internal/common"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/repo"
	"pg-sh-scripts/pkg/logging"

	uuid "github.com/satori/go.uuid"
)

type (
	IBashService interface {
		GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error)
		GetAll(ctx context.Context) ([]*model.Bash, error)
		CreateBash(ctx context.Context, dto dto.CreateBashDTO) (*model.Bash, error)
	}

	BashService struct {
		repository repo.IBashRepository
		logger     *logging.Logger
	}
)

func (s *BashService) GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error) {
	bash, err := s.repository.GetOneById(ctx, id)
	if err != nil {
		return nil, err
	}
	return bash, nil
}

func (s *BashService) GetAll(ctx context.Context) ([]*model.Bash, error) {
	bashList, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return bashList, nil
}

func (s *BashService) CreateBash(ctx context.Context, dto dto.CreateBashDTO) (*model.Bash, error) {
	bash, err := s.repository.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	return bash, nil
}

func GetBashService() IBashService {
	return &BashService{
		repository: repo.GetPgBashRepository(),
		logger:     common.GetLogger(),
	}
}
