package postgres

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primary_key"`
	Name     string
	User     uuid.UUID
	Tags     []*Tag `gorm:"many2many:album_tags;"`
	Pictures []*Picture
	IsPublic bool
}

type albumRepo struct {
	db *gorm.DB
}

func albumToDomain(album *Album) *domain.Album {
	return &domain.Album{
		ID:          album.ID,
		Name:        album.Name,
		User:        album.User,
		CreatedDate: album.CreatedAt,
		IsPublic:    album.IsPublic,
		Tags:        tagsToDomain(album.Tags),
	}
}

func albumFromDomain(album *domain.Album) *Album {
	var tags []*Tag

	for _, tag := range album.Tags {
		tags = append(tags, &Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return &Album{
		ID:       album.ID,
		Name:     album.Name,
		User:     album.User,
		IsPublic: album.IsPublic,
		Tags:     tags,
	}
}

func albumsToDomain(albums []*Album) []*domain.Album {
	albumList := []*domain.Album{}

	for _, album := range albums {
		albumList = append(albumList, albumToDomain(album))
	}

	return albumList
}

func (a albumRepo) CreateAlbum(album *domain.Album) error {
	albumToCreate := albumFromDomain(album)

	result := a.db.Create(albumToCreate)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a albumRepo) FindByName(user uuid.UUID, name string) (*domain.Album, error) {
	var album *Album

	query := a.db.Preload("Tags").Where("\"user\" = ? AND name = ?", user, name).First(&album)

	if query.Error != nil {
		return nil, query.Error
	}

	return albumToDomain(album), nil
}

func (a albumRepo) FindByUser(user uuid.UUID) ([]*domain.Album, error) {
	var albums []*Album

	query := a.db.Preload("Tags").Where("\"user\" = ?", user).Find(&albums)

	if query.Error != nil {
		return nil, query.Error
	}

	return albumsToDomain(albums), nil
}

func (a albumRepo) FindById(user uuid.UUID, id string) (*domain.Album, error) {
	var album *Album

	query := a.db.Preload("Tags").Where("\"user\" = ? AND \"id\" = ?", user, id).First(&album)

	if query.Error != nil {
		return nil, query.Error
	}

	return albumToDomain(album), nil
}

func (a albumRepo) DeleteAlbum(album *domain.Album) error {
	idToRemove := album.ID

	query := a.db.Where("id = ?", idToRemove).Delete(&Album{})

	return query.Error
}

func (a albumRepo) UpdateAlbum(album *domain.Album) error {
	albumToUpdate := albumFromDomain(album)

	err := a.db.Model(&albumToUpdate).Association("Tags").Replace(albumToUpdate.Tags)
	a.db.Omit("CreatedAt").Save(&albumToUpdate)

	return err
}

func NewAlbumRepo(db *gorm.DB) albumRepo {
	return albumRepo{
		db: db,
	}
}
