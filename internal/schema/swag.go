package schema

import "pg-sh-scripts/internal/model"

type (
	BashPaginationPage struct {
		Items  []model.Bash `json:"items"`
		Limit  int          `json:"limit"`
		Offset int          `json:"offset"`
		Total  int          `json:"total"`
	}

	BashLogPaginationPage struct {
		Items  []model.BashLog `json:"items"`
		Limit  int             `json:"limit"`
		Offset int             `json:"offset"`
		Total  int             `json:"total"`
	}
)
