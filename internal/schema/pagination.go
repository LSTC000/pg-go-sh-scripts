package schema

import "pg-sh-scripts/internal/model"

type (
	PaginationParams struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	PaginationLimitOffsetPage[T any] struct {
		Items  []T `json:"items"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	SwagBashPaginationLimitOffsetPage struct {
		Items  []model.Bash `json:"items"`
		Limit  int          `json:"limit"`
		Offset int          `json:"offset"`
		Total  int          `json:"total"`
	}

	SwagBashLogPaginationLimitOffsetPage struct {
		Items  []model.BashLog `json:"items"`
		Limit  int             `json:"limit"`
		Offset int             `json:"offset"`
		Total  int             `json:"total"`
	}
)
