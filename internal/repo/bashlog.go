package repo

import (
	"context"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/schema"

	uuid "github.com/satori/go.uuid"
)

type IBashLogRepository interface {
	GetPaginationPageByBashId(ctx context.Context, bashId uuid.UUID, paginationParams schema.PaginationParams) (schema.PaginationPage[*model.BashLog], error)
	Create(ctx context.Context, dto dto.CreateBashLogDTO) (*model.BashLog, error)
}
