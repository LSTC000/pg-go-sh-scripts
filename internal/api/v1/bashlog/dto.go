package bashlog

import uuid "github.com/satori/go.uuid"

type CreateBashLogDTO struct {
	BashID uuid.UUID `json:"bashId" swaggertype:"primitive,string"`
	Body   string    `json:"body"`
}
