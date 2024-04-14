package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Bash struct {
	Id        uuid.UUID `json:"id" swaggertype:"primitive,string" example:"59628b82-356c-4745-bc81-187015cde387"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt" example:"2024-04-14T15:50:21.907561+07:00"`
}
