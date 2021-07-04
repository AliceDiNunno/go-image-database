package domain

import (
	"github.com/google/uuid"
	"time"
)

type Album struct {
	ID   uuid.UUID
	User uuid.UUID

	CreatedDate time.Time

	Parent   *Album
	Name     string
	Tags     []*Tag
	IsPublic bool
}

func (a *Album) Initialize() {
	a.ID = uuid.New()
}
