package pagination

type (
	LimitOffsetParams struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}

	LimitOffsetPage[T any] struct {
		Items  []T `json:"items"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}
)
