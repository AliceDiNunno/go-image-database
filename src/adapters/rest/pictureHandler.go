package rest

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/AliceDiNunno/go-image-database/src/core/domain/Request"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func (rH RoutesHandler) fetchingPictureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := rH.getAuthenticatedUser(c)
		if user == nil {
			return
		}

		album := rH.getAlbum(c)
		if album == nil {
			return
		}

		id, stderr := uuid.Parse(c.Param("picture"))

		if stderr != nil {
			rH.handleError(c, e.Wrap(stderr).Append(ErrFormValidation))
			return
		}

		picture, err := rH.usecases.FetchPicture(user, album, id)

		if err != nil {
			rH.handleError(c, err.Append(domain.ErrPictureNotFound))
			return
		}

		c.Set("picture", picture)
	}
}

func (rH RoutesHandler) getPicture(c *gin.Context) *domain.Picture {
	picture, exists := c.Get("picture")

	if !exists {
		return nil
	}

	foundPicture := picture.(*domain.Picture)

	return foundPicture
}

func (rH RoutesHandler) GetPictureHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	album := rH.getAlbum(c)
	if album == nil {
		return
	}

	picture := rH.getPicture(c)
	if picture == nil {
		return
	}

	reader, err := rH.usecases.FetchPictureData(user, album, picture)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	io.Copy(c.Writer, reader)
}

//TODO: what happens if we upload a bigger file than the system could receive ? exploit ?
func (rH RoutesHandler) PostToAlbumHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	album := rH.getAlbum(c)
	if album == nil {
		return
	}

	file, header, stderr := c.Request.FormFile("upload")

	if stderr != nil {
		rH.handleError(c, e.Wrap(stderr).Append(ErrUploadFileNotFound))
		return
	}

	contentType := header.Header["Content-Type"]

	if len(contentType) <= 0 {
		rH.handleError(c, e.Wrap(ErrMissingContentType))
		return
	}

	err := rH.usecases.UploadPicture(user, album, file, contentType[0])

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (rH RoutesHandler) RemovePictureHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	album := rH.getAlbum(c)
	if album == nil {
		return
	}

	picture := rH.getPicture(c)
	if picture == nil {
		return
	}

	err := rH.usecases.DeletePicture(user, album, picture)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (rH RoutesHandler) EditPictureDataHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	album := rH.getAlbum(c)
	if album == nil {
		return
	}

	picture := rH.getPicture(c)
	if picture == nil {
		return
	}

	var request Request.EditPictureRequest
	stderr := c.ShouldBind(&request)

	if stderr != nil {
		rH.handleError(c, e.Wrap(ErrFormValidation))
		return
	}

	err := rH.usecases.UpdatePicture(user, album, picture, request)

	if err != nil {
		rH.handleError(c, err)
		return
	}
}
