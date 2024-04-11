package bash

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type IRepository interface {
	GetOneById(ctx context.Context, id uuid.UUID) (*Bash, error)
	GetAll(ctx context.Context) ([]*Bash, error)
	Create(ctx context.Context, createBash CreateBashDTO) (*Bash, error)
}
