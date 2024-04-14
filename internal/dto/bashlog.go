package dto

import uuid "github.com/satori/go.uuid"

type CreateBashLogDTO struct {
	BashId  uuid.UUID `json:"bashId" swaggertype:"primitive,string" example:"59628b82-356c-4745-bc81-187015cde387"`
	Body    string    `json:"body"`
	IsError bool      `json:"isError"`
}
