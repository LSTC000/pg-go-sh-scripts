package repo

import (
	"context"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/pkg/sql/pagination"

	uuid "github.com/satori/go.uuid"
)

type IBashRepository interface {
	GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error)
	GetPaginationPage(ctx context.Context, paginationParams pagination.LimitOffsetParams) (alias.BashLimitOffsetPage, error)
	Create(ctx context.Context, dto dto.CreateBash) (*model.Bash, error)
	RemoveById(ctx context.Context, id uuid.UUID) (*model.Bash, error)
}
