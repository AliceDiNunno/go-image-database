package usecases

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/google/uuid"
	"io"
)

type Logger interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

type AlbumRepo interface {
	CreateAlbum(album *domain.Album) *e.Error
	FindByName(user uuid.UUID, name string) (*domain.Album, *e.Error)
	FindByUser(user uuid.UUID) ([]*domain.Album, *e.Error)
	FindById(user uuid.UUID, id uuid.UUID) (*domain.Album, *e.Error)
	DeleteAlbum(album *domain.Album) *e.Error
	UpdateAlbum(album *domain.Album) *e.Error
}

type PictureRepo interface {
	SearchPictures(album *domain.Album, tags []string) ([]*domain.SearchPictureResult, *e.Error)
	FindPictures(album *domain.Album) ([]*domain.Picture, *e.Error)
	FindById(user uuid.UUID, album uuid.UUID, picture uuid.UUID) (*domain.Picture, *e.Error)
	CreatePicture(picture *domain.Picture) *e.Error
	DeletePicture(picture *domain.Picture) *e.Error
	UpdatePicture(picture *domain.Picture) *e.Error
}

type TagRepo interface {
	FindTag(name string) (*domain.Tag, *e.Error)
	FindTags(name string) ([]*domain.Tag, *e.Error)
	CreateTag(tag *domain.Tag) *e.Error
}

type FileStorage interface {
	UploadFile(id string, file io.Reader) *e.Error
	DeleteFile(id string) *e.Error
	GetFile(id string) (io.Reader, *e.Error)
}

type interactor struct {
	albumRepo   AlbumRepo
	pictureRepo PictureRepo
	tagRepo     TagRepo
	fileStorage FileStorage
}

func NewInteractor(aR AlbumRepo, pR PictureRepo, tR TagRepo, fS FileStorage) interactor {
	return interactor{
		albumRepo:   aR,
		pictureRepo: pR,
		tagRepo:     tR,
		fileStorage: fS,
	}
}
