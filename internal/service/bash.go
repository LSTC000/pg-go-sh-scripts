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
	IBashService interface {
		GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error)
		GetPaginationPage(ctx context.Context, paginationParams pagination.LimitOffsetParams) (alias.BashLimitOffsetPage, error)
		Create(ctx context.Context, dto dto.CreateBash) (*model.Bash, error)
		RemoveById(ctx context.Context, id uuid.UUID) (*model.Bash, error)
	}

	BashService struct {
		repository repo.IBashRepository
	}
)

func (s *BashService) GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error) {
	bash, err := s.repository.GetOneById(ctx, id)
	if err != nil {
		return nil, err
	}
	return bash, nil
}

func (s *BashService) GetPaginationPage(ctx context.Context, paginationParams pagination.LimitOffsetParams) (alias.BashLimitOffsetPage, error) {
	bashPaginationPage, err := s.repository.GetPaginationPage(ctx, paginationParams)
	if err != nil {
		return bashPaginationPage, err
	}
	return bashPaginationPage, nil
}

func (s *BashService) Create(ctx context.Context, dto dto.CreateBash) (*model.Bash, error) {
	bash, err := s.repository.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	return bash, nil
}

func (s *BashService) RemoveById(ctx context.Context, id uuid.UUID) (*model.Bash, error) {
	bash, err := s.repository.RemoveById(ctx, id)
	if err != nil {
		return nil, err
	}
	return bash, nil
}

func GetBashService() IBashService {
	return &BashService{
		repository: repo.GetPgBashRepository(),
	}
}
