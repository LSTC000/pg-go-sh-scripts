package repo

import (
	"context"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/pkg/sql/pagination"

	uuid "github.com/satori/go.uuid"
)

type IBashLogRepository interface {
	GetPaginationPageByBashId(ctx context.Context, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) (pagination.LimitOffsetPage[*model.BashLog], error)
	Create(ctx context.Context, dto dto.CreateBashLogDTO) (*model.BashLog, error)
}
