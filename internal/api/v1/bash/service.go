package bash

import (
	"context"
	"pg-sh-scripts/internal/common"
	"pg-sh-scripts/pkg/logging"

	uuid "github.com/satori/go.uuid"
)

type (
	IService interface {
		GetOneById(ctx context.Context, id uuid.UUID) (*Bash, error)
		GetAll(ctx context.Context) ([]*Bash, error)
		CreateBash(ctx context.Context, createBash CreateBashDTO) (*Bash, error)
	}

	Service struct {
		repository IRepository
		logger     *logging.Logger
	}
)

func (s *Service) GetOneById(ctx context.Context, id uuid.UUID) (*Bash, error) {
	bash, err := s.repository.GetOneById(ctx, id)
	if err != nil {
		return nil, err
	}
	return bash, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*Bash, error) {
	bashList, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return bashList, nil
}

func (s *Service) CreateBash(ctx context.Context, createBash CreateBashDTO) (*Bash, error) {
	bash, err := s.repository.Create(ctx, createBash)
	if err != nil {
		return nil, err
	}
	return bash, nil
}

func GetService() IService {
	return &Service{
		repository: GetPgRepository(),
		logger:     common.GetLogger(),
	}
}
