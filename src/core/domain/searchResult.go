package domain

import "github.com/google/uuid"

type SearchPictureResult struct {
	ID      uuid.UUID
	AlbumID uuid.UUID
}
