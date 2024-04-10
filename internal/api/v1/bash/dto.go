package bash

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type (
	CreateBashDTO struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	ExecBashDTO struct {
		ID             uuid.UUID     `json:"id" swaggertype:"primitive,string"`
		TimeoutSeconds time.Duration `json:"timeoutSeconds" swaggertype:"primitive,integer"`
	}
)
