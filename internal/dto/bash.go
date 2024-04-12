package dto

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
		Id             uuid.UUID     `json:"id" swaggertype:"primitive,string"`
		TimeoutSeconds time.Duration `json:"timeoutSeconds" swaggertype:"primitive,integer"`
	}
)