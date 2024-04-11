package bashlog

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

type IRepository interface {
	GetAllByBashID(ctx context.Context, bashID uuid.UUID) ([]*BashLog, error)
	Create(ctx context.Context, createBashLog CreateBashLogDTO) (*BashLog, error)
}
