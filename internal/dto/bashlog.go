package dto

import uuid "github.com/satori/go.uuid"

type CreateBashLogDTO struct {
	BashId uuid.UUID `json:"bashId" swaggertype:"primitive,string"`
	Body   string    `json:"body"`
}
