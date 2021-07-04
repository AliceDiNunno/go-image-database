package usecases

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/AliceDiNunno/go-image-database/src/core/domain/Request"
	"io"
)

type Usecases interface {
	CreateAlbum(user *domain.User, request *Request.CreateAlbumRequest) error
	GetUserAlbums(user *domain.User) ([]*domain.Album, error)
	DeleteAlbum(user *domain.User, id string) error
	GetAlbumsContent(user *domain.User, id string) ([]*domain.Picture, error)
	UpdateAlbum(user *domain.User, albumId string, request Request.EditAlbumRequest) error

	UploadPicture(user *domain.User, albumId string, file io.Reader, contentType string) error
	DeletePicture(user *domain.User, albumId string, fileId string) error
	FetchPicture(user *domain.User, albumId string, fileId string) (io.Reader, error)
	UpdatePicture(user *domain.User, albumId string, fileId string, request Request.EditPictureRequest) error
}
