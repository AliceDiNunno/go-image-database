package rest

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (rH RoutesHandler) GetPictureHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	if user == nil {
		return
	}

	albumId := c.Param("album")
	pictureId := c.Param("picture")

	reader, err := rH.usecases.FetchPicture(user, albumId, pictureId)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	io.Copy(c.Writer, reader)
}

//TODO: what happens if we upload a bigger file than the system could receive ? exploit ?
func (rH RoutesHandler) PostToAlbumHandler(c *gin.Context) {
	user := rH.getAuthenticatedUser(c)
	id := c.Param("album")
	file, header, err := c.Request.FormFile("upload")

	if err != nil {
		rH.handleError(c, ErrUploadFileNotFound)
		return
	}

	contentType := header.Header["Content-Type"]

	if len(contentType) <= 0 {
		rH.handleError(c, ErrFormValidation)
		return
	}

	if err != nil {
		rH.handleError(c, ErrFormValidation)
		return
	}

	err = rH.usecases.UploadPicture(user, id, file, contentType[0])

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

	albumId := c.Param("album")
	pictureId := c.Param("picture")

	err := rH.usecases.DeletePicture(user, albumId, pictureId)

	if err != nil {
		rH.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (rH RoutesHandler) EditPictureDataHandler(c *gin.Context) {

}
