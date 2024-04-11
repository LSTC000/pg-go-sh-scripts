package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Bash struct {
	Id        uuid.UUID `json:"id" swaggertype:"primitive,string"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}
