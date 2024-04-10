package bash

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Bash struct {
	ID        uuid.UUID `json:"id" swaggertype:"primitive,string"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}
