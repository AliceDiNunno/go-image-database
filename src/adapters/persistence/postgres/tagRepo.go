package postgres

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primary_key"`
	Name string    `gorm:"uniqueIndex"`
}

type tagRepo struct {
	db *gorm.DB
}

func tagToDomain(tag *Tag) *domain.Tag {
	return &domain.Tag{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func tagsToDomain(tags []*Tag) []*domain.Tag {
	tagList := []*domain.Tag{}

	for _, tag := range tags {
		tagList = append(tagList, tagToDomain(tag))
	}

	return tagList
}

func tagFromDomain(tag *domain.Tag) *Tag {
	return &Tag{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func (t tagRepo) CreateTag(tag *domain.Tag) *e.Error {
	tagToCreate := tagFromDomain(tag)

	result := t.db.Create(tagToCreate)

	if result.Error != nil {
		return e.Wrap(result.Error)
	}

	return nil
}

func (t tagRepo) DeleteTag(tag *domain.Tag) *e.Error {
	idToRemove := tag.ID

	query := t.db.Where("id = ?", idToRemove).Delete(&Tag{})

	return e.Wrap(query.Error)
}

func (t tagRepo) FindTag(name string) (*domain.Tag, *e.Error) {
	var tag *domain.Tag

	query := t.db.Where("name = ?", name).First(&tag)

	if query.Error != nil {
		return nil, e.Wrap(query.Error)
	}

	return tag, nil
}

func (t tagRepo) FindTags(name string) ([]*domain.Tag, *e.Error) {
	var tags []*domain.Tag

	query := t.db.Where("name = ?", name).Find(&tags)

	if query.Error != nil {
		return nil, e.Wrap(query.Error)
	}

	return tags, nil
}

func NewTagRepo(db *gorm.DB) tagRepo {
	return tagRepo{
		db: db,
	}
}
