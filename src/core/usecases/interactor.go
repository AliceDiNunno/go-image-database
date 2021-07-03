package usecases

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
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
	CreateAlbum(album *domain.Album) error
	FindByName(user uuid.UUID, name string) (*domain.Album, error)
	FindByUser(user uuid.UUID) ([]*domain.Album, error)
	FindById(user uuid.UUID, id string) (*domain.Album, error)
	DeleteAlbum(album *domain.Album) error
}

type PictureRepo interface {
	CreatePicture(picture *domain.Picture) error
	DeletePicture(picture *domain.Picture) error
}

type FileStorage interface {
	UploadFile(id string, file io.Reader) error
	DeleteFile(id string) error
}

type interactor struct {
	albumRepo   AlbumRepo
	pictureRepo PictureRepo
	fileStorage FileStorage
}

func NewInteractor(aR AlbumRepo, pR PictureRepo, fS FileStorage) interactor {
	return interactor{
		albumRepo:   aR,
		pictureRepo: pR,
		fileStorage: fS,
	}
}
