package bashlog

import (
	"context"
	"pg-sh-scripts/internal/common"
	"pg-sh-scripts/pkg/logging"

	uuid "github.com/satori/go.uuid"
)

type (
	IService interface {
		GetAllByBashId(ctx context.Context, bashId uuid.UUID) ([]*BashLog, error)
	}

	Service struct {
		repository IRepository
		logger     *logging.Logger
	}
)

func (s *Service) GetAllByBashId(ctx context.Context, bashId uuid.UUID) ([]*BashLog, error) {
	bashLogList, err := s.repository.GetAllByBashId(ctx, bashId)
	if err != nil {
		return nil, err
	}
	return bashLogList, nil
}

func GetService() IService {
	return &Service{
		repository: GetPgRepository(),
		logger:     common.GetLogger(),
	}
}
