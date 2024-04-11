package repo

import (
	"context"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"

	uuid "github.com/satori/go.uuid"
)

type IBashRepository interface {
	GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error)
	GetAll(ctx context.Context) ([]*model.Bash, error)
	Create(ctx context.Context, dto dto.CreateBashDTO) (*model.Bash, error)
}
