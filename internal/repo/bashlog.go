package repo

import (
	"context"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/pkg/sql/pagination"

	uuid "github.com/satori/go.uuid"
)

type IBashLogRepository interface {
	GetPaginationPageByBashId(ctx context.Context, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) (alias.BashLogLimitOffsetPage, error)
	Create(ctx context.Context, dto dto.CreateBashLog) (*model.BashLog, error)
}
