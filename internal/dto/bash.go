package dto

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type (
	CreateBash struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	ExecBash struct {
		Id             uuid.UUID     `json:"id" swaggertype:"primitive,string" example:"59628b82-356c-4745-bc81-187015cde387"`
		TimeoutSeconds time.Duration `json:"timeoutSeconds" swaggertype:"primitive,integer"`
	}
)
