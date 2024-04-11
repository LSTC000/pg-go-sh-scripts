package repo

import (
	"context"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"

	uuid "github.com/satori/go.uuid"
)

type IBashLogRepository interface {
	GetAllByBashId(ctx context.Context, bashId uuid.UUID) ([]*model.BashLog, error)
	Create(ctx context.Context, dto dto.CreateBashLogDTO) (*model.BashLog, error)
}
