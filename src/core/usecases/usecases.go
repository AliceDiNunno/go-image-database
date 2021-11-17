package usecases

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/AliceDiNunno/go-image-database/src/core/domain/Request"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/google/uuid"
	"io"
)

type Usecases interface {
	CreateAlbum(user *domain.User, request *Request.CreateAlbumRequest) *e.Error
	GetUserAlbums(user *domain.User) ([]*domain.Album, *e.Error)
	FetchAlbum(user *domain.User, id uuid.UUID) (*domain.Album, *e.Error)
	DeleteAlbum(user *domain.User, album *domain.Album) *e.Error
	GetAlbumsContent(user *domain.User, album *domain.Album) ([]*domain.Picture, *e.Error)
	UpdateAlbum(user *domain.User, album *domain.Album, request Request.EditAlbumRequest) *e.Error
	SearchAlbumContent(user *domain.User, album *domain.Album, request Request.SearchAlbumRequest) ([]*domain.SearchPictureResult, *e.Error)

	FetchPicture(user *domain.User, album *domain.Album, pictureId uuid.UUID) (*domain.Picture, *e.Error)
	FetchPictureData(user *domain.User, album *domain.Album, picture *domain.Picture) (io.Reader, *e.Error)
	DeletePicture(user *domain.User, album *domain.Album, picture *domain.Picture) *e.Error
	UpdatePicture(user *domain.User, album *domain.Album, picture *domain.Picture, request Request.EditPictureRequest) *e.Error
	UploadPicture(user *domain.User, album *domain.Album, picture io.Reader, contentType string) *e.Error
}
