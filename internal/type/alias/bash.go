package alias

import (
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/pkg/sql/pagination"
)

type (
	BashTitle           = string
	BashLimitOffsetPage = pagination.LimitOffsetPage[*model.Bash]
)
