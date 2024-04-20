package repo

import (
	"context"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/pkg/sql/pagination"

	uuid "github.com/satori/go.uuid"
)

type IBashRepository interface {
	GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error)
	GetPaginationPage(ctx context.Context, paginationParams pagination.LimitOffsetParams) (pagination.LimitOffsetPage[*model.Bash], error)
	Create(ctx context.Context, dto dto.CreateBashDTO) (*model.Bash, error)
	RemoveById(ctx context.Context, id uuid.UUID) (*model.Bash, error)
}
