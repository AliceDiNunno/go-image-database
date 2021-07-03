package postgres

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Picture struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primary_key"`
	User uuid.UUID `gorm:"type:uuid"`

	CreatedDate time.Time

	Tags    []Tag `gorm:"many2many:picture_tags;"`
	AlbumID uuid.UUID
	Album   *Album
}

type pictureRepo struct {
	db *gorm.DB
}

func pictureToDomain(picture *Picture) *domain.Picture {
	return &domain.Picture{
		ID:          picture.ID,
		User:        picture.User,
		CreatedDate: picture.CreatedAt,
		Album:       albumToDomain(picture.Album),
		Tags:        nil,
	}
}

func picturesToDomain(pictures []*Picture) []*domain.Picture {
	pictureList := []*domain.Picture{}

	for _, picture := range pictures {
		pictureList = append(pictureList, pictureToDomain(picture))
	}

	return pictureList
}

func pictureFromDomain(picture *domain.Picture) *Picture {
	return &Picture{
		ID:          picture.ID,
		User:        picture.User,
		CreatedDate: picture.CreatedDate,
		AlbumID:     picture.Album.ID,
		Tags:        nil,
	}
}

func (p pictureRepo) CreatePicture(picture *domain.Picture) error {
	pictureToCreate := pictureFromDomain(picture)

	result := p.db.Create(pictureToCreate)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (p pictureRepo) DeletePicture(picture *domain.Picture) error {
	idToRemove := picture.ID

	query := p.db.Where("id = ?", idToRemove).Delete(&Picture{})

	return query.Error
}

func NewPictureRepo(db *gorm.DB) pictureRepo {
	return pictureRepo{
		db: db,
	}
}
