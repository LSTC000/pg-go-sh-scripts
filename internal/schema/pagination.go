package schema

import "pg-sh-scripts/internal/model"

type (
	PaginationParams struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	PaginationPage[T any] struct {
		Items  []T `json:"items"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	SwagBashPaginationPage struct {
		Items  []model.Bash `json:"items"`
		Limit  int          `json:"limit"`
		Offset int          `json:"offset"`
		Total  int          `json:"total"`
	}

	SwagBashLogPaginationPage struct {
		Items  []model.BashLog `json:"items"`
		Limit  int             `json:"limit"`
		Offset int             `json:"offset"`
		Total  int             `json:"total"`
	}
)
