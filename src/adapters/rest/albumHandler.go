package rest

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain/Request"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
	err := c.ShouldBind(&request)
	if err != nil {
		rH.handleError(c, ErrFormValidation)
		return
	}

	err = rH.usecases.CreateAlbum(user, &request)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (rH RoutesHandler) GetAlbumContentHandler(c *gin.Context) {
	id := c.Param("album")
	user := rH.getAuthenticatedUser(c)

	content, err := rH.usecases.GetAlbumsContent(user, id)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, content)
}

func (rH RoutesHandler) FindAlbumContentHandler(c *gin.Context) {

}

func (rH RoutesHandler) DeleteAlbumHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	id := c.Param("album")

	err := rH.usecases.DeleteAlbum(user, id)

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
	var request Request.EditAlbumRequest
	err := c.ShouldBind(&request)

	if err != nil {
		rH.handleError(c, ErrFormValidation)
		return
	}

	albumId := c.Param("album")

	err = rH.usecases.UpdateAlbum(user, albumId, request)

	if err != nil {
		rH.handleError(c, err)
		return
	}
}
