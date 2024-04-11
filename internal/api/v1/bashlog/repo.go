package bashlog

import (
	"context"
)

type IRepository interface {
	Create(ctx context.Context, createBashLog CreateBashLogDTO) (*BashLog, error)
}
