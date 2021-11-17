package postgres

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	e "github.com/AliceDiNunno/go-nested-traced-error"
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

func (a albumRepo) CreateAlbum(album *domain.Album) *e.Error {
	albumToCreate := albumFromDomain(album)

	result := a.db.Create(albumToCreate)

	if result.Error != nil {
		return e.Wrap(result.Error)
	}

	return nil
}

func (a albumRepo) FindByName(user uuid.UUID, name string) (*domain.Album, *e.Error) {
	var album *Album

	query := a.db.Preload("Tags").Where("\"user\" = ? AND name = ?", user, name).First(&album)

	if query.Error != nil {
		return nil, e.Wrap(query.Error)
	}

	return albumToDomain(album), nil
}

func (a albumRepo) FindByUser(user uuid.UUID) ([]*domain.Album, *e.Error) {
	var albums []*Album

	query := a.db.Preload("Tags").Where("\"user\" = ?", user).Find(&albums)

	if query.Error != nil {
		return nil, e.Wrap(query.Error)
	}

	return albumsToDomain(albums), nil
}

func (a albumRepo) FindById(user uuid.UUID, id uuid.UUID) (*domain.Album, *e.Error) {
	var album *Album

	query := a.db.Preload("Tags").Where("(\"user\" = ? OR \"is_public\" = true) AND \"id\" = ?", user, id).First(&album)

	if query.Error != nil {
		return nil, e.Wrap(query.Error)
	}

	return albumToDomain(album), nil
}

func (a albumRepo) DeleteAlbum(album *domain.Album) *e.Error {
	idToRemove := album.ID

	query := a.db.Where("id = ?", idToRemove).Delete(&Album{})

	return e.Wrap(query.Error)
}

func (a albumRepo) UpdateAlbum(album *domain.Album) *e.Error {
	albumToUpdate := albumFromDomain(album)

	err := a.db.Model(&albumToUpdate).Association("Tags").Replace(albumToUpdate.Tags)
	a.db.Omit("CreatedAt").Save(&albumToUpdate)

	return e.Wrap(err)
}

func NewAlbumRepo(db *gorm.DB) albumRepo {
	return albumRepo{
		db: db,
	}
}
