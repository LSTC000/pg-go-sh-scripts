package bashlog

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type BashLog struct {
	Id        uuid.UUID `json:"id" swaggertype:"primitive,string"`
	BashId    uuid.UUID `json:"bashId" swaggertype:"primitive,string"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}
