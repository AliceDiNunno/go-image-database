package postgres

import (
	"fmt"
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Picture struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primary_key"`
	User uuid.UUID `gorm:"type:uuid"`

	Tags    []*Tag `gorm:"many2many:picture_tags;"`
	AlbumID uuid.UUID
	Album   *Album
	Phash   string
}

type PictureSearchResult struct {
	ID      uuid.UUID
	AlbumID uuid.UUID
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
		Tags:        tagsToDomain(picture.Tags),
		Phash:       picture.Phash,
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
	var tags []*Tag

	for _, tag := range picture.Tags {
		tags = append(tags, &Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return &Picture{
		ID:      picture.ID,
		User:    picture.User,
		AlbumID: picture.Album.ID,
		Tags:    tags,
		Phash:   picture.Phash,
	}
}

func searchPictureToDomain(picture *PictureSearchResult) *domain.SearchPictureResult {
	return &domain.SearchPictureResult{
		ID:      picture.ID,
		AlbumID: picture.AlbumID,
	}
}

func searchPicturesToDomain(pictures []*PictureSearchResult) []*domain.SearchPictureResult {
	albumList := []*domain.SearchPictureResult{}

	for _, album := range pictures {
		albumList = append(albumList, searchPictureToDomain(album))
	}

	return albumList
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

func (p pictureRepo) FindPictures(album *domain.Album) ([]*domain.Picture, error) {
	albumId := album.ID
	var pictures []*Picture

	query := p.db.Joins("Album").Preload("Tags").Where("album_id = ?", albumId).Find(&pictures)

	if query.Error != nil {
		return nil, query.Error
	}

	return picturesToDomain(pictures), nil
}

func (p pictureRepo) FindById(userId uuid.UUID, albumId uuid.UUID, pictureId uuid.UUID) (*domain.Picture, error) {
	var picture *Picture

	query := p.db.Joins("Album").Preload("Tags").Where("pictures.album_id = ? AND pictures.id = ?", albumId, pictureId).Find(&picture)

	if query.Error != nil {
		return nil, query.Error
	}

	if query.RowsAffected <= 0 {
		return nil, domain.ErrPictureNotFound
	}

	if !picture.Album.IsPublic && picture.User != userId {
		return nil, domain.ErrPictureNotFound
	}

	return pictureToDomain(picture), nil
}

func (p pictureRepo) UpdatePicture(picture *domain.Picture) error {
	pictureToUpdate := pictureFromDomain(picture)

	err := p.db.Model(&pictureToUpdate).Association("Tags").Replace(pictureToUpdate.Tags)

	return err
}

//TODO: there should be 1000 ways to improve this but I guess I don't master gorm well enough yet
func (p pictureRepo) SearchPictures(album *domain.Album, tags []string) ([]*domain.SearchPictureResult, error) {
	whereClause := "WHERE "

	if len(tags) > 0 {
		for idx, tag := range tags {
			if idx > 0 {
				whereClause = fmt.Sprintf("%s AND ", whereClause)
			}

			whereClause = fmt.Sprintf("%s '%s' = ANY(tag_array)", whereClause, tag)
		}
	}

	var pictures []*PictureSearchResult

	rawQuery := "SELECT p.id, p.created_at, p.updated_at, p.deleted_at, p.user, p.album_id, t.tag_array FROM pictures p, LATERAL ( " +
		"SELECT ARRAY ( " +
		"SELECT pt.tag_id " +
		"FROM picture_tags pt " +
		"WHERE pt.picture_id = p.id " +
		") AS tag_array " +
		") t"

	rawQuery = fmt.Sprintf("%s %s ;", rawQuery, whereClause)

	query := p.db.Raw(rawQuery).Scan(&pictures)

	return searchPicturesToDomain(pictures), query.Error
}

func NewPictureRepo(db *gorm.DB) pictureRepo {
	return pictureRepo{
		db: db,
	}
}
