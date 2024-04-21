package alias

import (
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/pkg/sql/pagination"
)

type BashLogLimitOffsetPage = pagination.LimitOffsetPage[*model.BashLog]
