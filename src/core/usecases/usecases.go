package usecases

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/AliceDiNunno/go-image-database/src/core/domain/Request"
	"github.com/google/uuid"
	"io"
)

type Usecases interface {
	CreateAlbum(user *domain.User, request *Request.CreateAlbumRequest) error
	GetUserAlbums(user *domain.User) ([]*domain.Album, error)
	FetchAlbum(user *domain.User, id uuid.UUID) (*domain.Album, error)
	DeleteAlbum(user *domain.User, album *domain.Album) error
	GetAlbumsContent(user *domain.User, album *domain.Album) ([]*domain.Picture, error)
	UpdateAlbum(user *domain.User, album *domain.Album, request Request.EditAlbumRequest) error
	SearchAlbumContent(user *domain.User, album *domain.Album, request Request.SearchAlbumRequest) ([]*domain.SearchPictureResult, error)

	FetchPicture(user *domain.User, album *domain.Album, pictureId uuid.UUID) (*domain.Picture, error)
	FetchPictureData(user *domain.User, album *domain.Album, picture *domain.Picture) (io.Reader, error)
	DeletePicture(user *domain.User, album *domain.Album, picture *domain.Picture) error
	UpdatePicture(user *domain.User, album *domain.Album, picture *domain.Picture, request Request.EditPictureRequest) error
	UploadPicture(user *domain.User, album *domain.Album, picture io.Reader, contentType string) error
}
