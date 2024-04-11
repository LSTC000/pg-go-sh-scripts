package bashlog

import uuid "github.com/satori/go.uuid"

type CreateBashLogDTO struct {
	BashID uuid.UUID `json:"bashID" swaggertype:"primitive,string"`
	Body   string    `json:"body"`
}
