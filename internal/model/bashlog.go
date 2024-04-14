package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type BashLog struct {
	Id        uuid.UUID `json:"id" swaggertype:"primitive,string" example:"f4f4d096-ef4a-4649-8346-a952e2ca27d3"`
	BashId    uuid.UUID `json:"bashId" swaggertype:"primitive,string" example:"59628b82-356c-4745-bc81-187015cde387"`
	Body      string    `json:"body"`
	IsError   bool      `json:"isError"`
	CreatedAt time.Time `json:"createdAt" example:"2024-04-14T15:50:21.907561+07:00"`
}
