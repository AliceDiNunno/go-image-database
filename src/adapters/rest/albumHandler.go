package rest

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/AliceDiNunno/go-image-database/src/core/domain/Request"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (rH RoutesHandler) fetchingAlbumMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := rH.getAuthenticatedUser(c)
		if user == nil {
			return
		}

		id, stderr := uuid.Parse(c.Param("album"))

		if stderr != nil {
			rH.handleError(c, e.Wrap(stderr).Append(ErrFormValidation))
			return
		}

		album, err := rH.usecases.FetchAlbum(user, id)

		if err != nil {
			rH.handleError(c, err.Append(domain.ErrAlbumNotFound))
			return
		}

		c.Set("album", album)
	}
}

func (rH RoutesHandler) getAlbum(c *gin.Context) *domain.Album {
	album, exists := c.Get("album")

	if !exists {
		return nil
	}

	foundAlbum := album.(*domain.Album)

	return foundAlbum
}

func (rH RoutesHandler) GetUserAlbumsHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	albums, err := rH.usecases.GetUserAlbums(user)

	if err != nil {
		rH.handleError(c, err)
	}

	c.JSON(http.StatusOK, albums)
}

func (rH RoutesHandler) CreateUserAlbumHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)

	var request Request.CreateAlbumRequest
	stderr := c.ShouldBind(&request)
	if stderr != nil {
		rH.handleError(c, e.Wrap(stderr).Append(ErrFormValidation))
		return
	}

	err := rH.usecases.CreateAlbum(user, &request)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (rH RoutesHandler) GetAlbumContentHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	album := rH.getAlbum(c)
	if album == nil {
		return
	}

	content, err := rH.usecases.GetAlbumsContent(user, album)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, content)
}

func (rH RoutesHandler) SearchAlbumContentHandler(c *gin.Context) {
	album := rH.getAlbum(c)
	if album == nil {
		return
	}

	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	var request Request.SearchAlbumRequest
	stderr := c.ShouldBind(&request)

	if stderr != nil {
		rH.handleError(c, e.Wrap(stderr).Append(ErrFormValidation))
		return
	}

	content, err := rH.usecases.SearchAlbumContent(user, album, request)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, content)
}

func (rH RoutesHandler) DeleteAlbumHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	album := rH.getAlbum(c)
	if album == nil {
		return
	}

	err := rH.usecases.DeleteAlbum(user, album)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (rH RoutesHandler) EditAlbumDataHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	album := rH.getAlbum(c)
	if album == nil {
		return
	}

	var request Request.EditAlbumRequest
	stderr := c.ShouldBind(&request)

	if stderr != nil {
		rH.handleError(c, e.Wrap(stderr).Append(ErrFormValidation))
		return
	}

	err := rH.usecases.UpdateAlbum(user, album, request)

	if err != nil {
		rH.handleError(c, err)
		return
	}
}
