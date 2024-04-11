package bashlog

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type IRepository interface {
	GetAllByBashId(ctx context.Context, bashId uuid.UUID) ([]*BashLog, error)
	Create(ctx context.Context, createBashLog CreateBashLogDTO) (*BashLog, error)
}
