package domain

import (
	"github.com/google/uuid"
	"time"
)

type Picture struct {
	ID   uuid.UUID
	User uuid.UUID

	CreatedDate time.Time

	Tags  []*Tag
	Album *Album
}

func (p *Picture) Initialize() {
	p.ID = uuid.New()
}
